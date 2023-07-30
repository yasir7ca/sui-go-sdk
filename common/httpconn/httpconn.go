package httpconn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/yasir7ca/sui-go-sdk/models"
)

type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (err HTTPError) Error() string {
	if len(err.Body) == 0 {
		return err.Status
	}
	return fmt.Sprintf("%v: %s", err.Status, err.Body)
}

const (
	vsn = "2.0"
)

var (
	ErrNoResult = errors.New("no result in JSON-RPC response")
)

type HttpConn struct {
	idCounter uint32
	rpcUrl    string
	client    *http.Client
}

func Dial(rpcUrl string) *HttpConn {
	hc := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    3,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
	return DialWithClient(rpcUrl, hc)
}

func DialWithClient(rpcUrl string, c *http.Client) *HttpConn {
	return &HttpConn{
		rpcUrl: strings.TrimRight(rpcUrl, "/"),
		client: c,
	}
}

// CallContext performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (h *HttpConn) CallContext(ctx context.Context, result interface{}, op Operation) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}
	msg, err := h.newMessage(op.Method, op.Params)
	if err != nil {
		return err
	}
	respBody, err := h.doRequest(ctx, msg)
	if err != nil {
		return err
	}
	defer respBody.Close()

	var respMsg models.JsonRPCMessage
	if err := json.NewDecoder(respBody).Decode(&respMsg); err != nil {
		return err
	}
	if respMsg.Error != nil {
		return respMsg.Error
	}
	if len(respMsg.Result) == 0 {
		return ErrNoResult
	}
	return json.Unmarshal(respMsg.Result, &result)
}

func (h *HttpConn) newMessage(method string, paramsIn ...interface{}) (*models.JsonRPCMessage, error) {
	msg := &models.JsonRPCMessage{Version: vsn, ID: h.nextID(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (h *HttpConn) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.rpcUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	req.Header.Set("Content-Type", "application/json")

	// do request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var body []byte
		if _, err := buf.ReadFrom(resp.Body); err == nil {
			body = buf.Bytes()
		}

		return nil, HTTPError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}
	return resp.Body, nil
}

func (h *HttpConn) nextID() json.RawMessage {
	id := atomic.AddUint32(&h.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

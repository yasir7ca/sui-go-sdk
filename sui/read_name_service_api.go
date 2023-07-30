package sui

import (
	"context"
	"errors"

	"github.com/yasir7ca/sui-go-sdk/common/httpconn"
	"github.com/yasir7ca/sui-go-sdk/models"
)

type IReadNameServiceFromSuiAPI interface {
	SuiXResolveNameServiceAddress(ctx context.Context, req models.SuiXResolveNameServiceAddressRequest) (string, error)
	SuiXResolveNameServiceNames(ctx context.Context, req models.SuiXResolveNameServiceNamesRequest) (models.SuiXResolveNameServiceNamesResponse, error)
}

type suiReadNameServiceFromSuiImpl struct {
	conn *httpconn.HttpConn
}

// SuiXResolveNameServiceAddress implements the method `suix_resolveNameServiceAddress`, get the resolved address given resolver and name.
func (s *suiReadNameServiceFromSuiImpl) SuiXResolveNameServiceAddress(ctx context.Context, req models.SuiXResolveNameServiceAddressRequest) (string, error) {
	return "", errors.New("not implemented")
}

// SuiXResolveNameServiceNames implements the method `suix_resolveNameServiceNames`, return the resolved names given address, if multiple names are resolved, the first one is the primary name.
func (s *suiReadNameServiceFromSuiImpl) SuiXResolveNameServiceNames(ctx context.Context, req models.SuiXResolveNameServiceNamesRequest) (models.SuiXResolveNameServiceNamesResponse, error) {
	var rsp models.SuiXResolveNameServiceNamesResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_resolveNameServiceNames",
		Params: []interface{}{
			req.Address,
			req.Cursor,
			req.Limit,
		},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

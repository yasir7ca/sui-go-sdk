// Copyright (c) BlockVision, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sui

import (
	"context"

	"github.com/yasir7ca/sui-go-sdk/common/httpconn"
)

type IBaseAPI interface {
	SuiCall(ctx context.Context, method string, params ...interface{}) (interface{}, error)
}

type suiBaseImpl struct {
	conn *httpconn.HttpConn
}

// SuiCall send customized request to Sui Node endpoint.
func (s *suiBaseImpl) SuiCall(ctx context.Context, method string, params ...interface{}) (interface{}, error) {
	var resp interface{}
	err := s.conn.CallContext(ctx, &resp, httpconn.Operation{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

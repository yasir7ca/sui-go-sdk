// Copyright (c) BlockVision, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sui

import (
	"context"

	"github.com/yasir7ca/sui-go-sdk/common/httpconn"
	"github.com/yasir7ca/sui-go-sdk/models"
)

type IReadEventFromSuiAPI interface {
	SuiGetEvents(ctx context.Context, req models.SuiGetEventsRequest) (models.GetEventsResponse, error)
	SuiXQueryEvents(ctx context.Context, req models.SuiXQueryEventsRequest) (models.PaginatedEventsResponse, error)
}

type suiReadEventFromSuiImpl struct {
	conn *httpconn.HttpConn
}

// SuiGetEvents implements the method `sui_getEvents`, gets transaction events.
func (s *suiReadEventFromSuiImpl) SuiGetEvents(ctx context.Context, req models.SuiGetEventsRequest) (models.GetEventsResponse, error) {
	var rsp models.GetEventsResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "sui_getEvents",
		Params: []interface{}{
			req.Digest,
		},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// SuiXQueryEvents implements the method `suix_queryEvents`, gets list of events for a specified query criteria.
func (s *suiReadEventFromSuiImpl) SuiXQueryEvents(ctx context.Context, req models.SuiXQueryEventsRequest) (models.PaginatedEventsResponse, error) {
	var rsp models.PaginatedEventsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_queryEvents",
		Params: []interface{}{
			req.SuiEventFilter,
			req.Cursor,
			req.Limit,
			req.DescendingOrder,
		},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

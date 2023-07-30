// Copyright (c) BlockVision, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sui

import (
	"context"

	"github.com/yasir7ca/sui-go-sdk/common/httpconn"
	"github.com/yasir7ca/sui-go-sdk/models"
)

type IReadSystemFromSuiAPI interface {
	SuiGetCheckpoint(ctx context.Context, req models.SuiGetCheckpointRequest) (models.CheckpointResponse, error)
	SuiGetCheckpoints(ctx context.Context, req models.SuiGetCheckpointsRequest) (models.PaginatedCheckpointsResponse, error)
	SuiGetLatestCheckpointSequenceNumber(ctx context.Context) (uint64, error)
	SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error)
	SuiXGetCommitteeInfo(ctx context.Context, req models.SuiXGetCommitteeInfoRequest) (models.SuiXGetCommitteeInfoResponse, error)
	SuiXGetStakes(ctx context.Context, req models.SuiXGetStakesRequest) ([]*models.DelegatedStakesResponse, error)
	SuiXGetStakesByIds(ctx context.Context, req models.SuiXGetStakesByIdsRequest) ([]*models.DelegatedStakesResponse, error)
	SuiXGetEpochs(ctx context.Context, req models.SuiXGetEpochsRequest) (models.PaginatedEpochInfoResponse, error)
	SuiXGetCurrentEpoch(ctx context.Context) (models.EpochInfo, error)
	SuiXGetLatestSuiSystemState(ctx context.Context) (models.SuiSystemStateSummary, error)
	SuiGetChainIdentifier(ctx context.Context) (string, error)
	SuiXGetValidatorsApy(ctx context.Context) (models.ValidatorsApy, error)
}

type suiReadSystemFromSuiImpl struct {
	conn *httpconn.HttpConn
}

// SuiGetCheckpoint implements the method `sui_getCheckpoint`, gets a checkpoint.
func (s *suiReadSystemFromSuiImpl) SuiGetCheckpoint(ctx context.Context, req models.SuiGetCheckpointRequest) (models.CheckpointResponse, error) {
	var rsp models.CheckpointResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "sui_getCheckpoint",
		Params: []interface{}{
			req.CheckpointID,
		},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiGetCheckpoints implements the method `sui_getCheckpoints`, gets paginated list of checkpoints.
func (s *suiReadSystemFromSuiImpl) SuiGetCheckpoints(ctx context.Context, req models.SuiGetCheckpointsRequest) (models.PaginatedCheckpointsResponse, error) {
	var rsp models.PaginatedCheckpointsResponse
	if err := validate.ValidateStruct(req); err != nil {
		return rsp, err
	}
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "sui_getCheckpoints",
		Params: []interface{}{
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

// SuiGetLatestCheckpointSequenceNumber implements the method `sui_getLatestCheckpointSequenceNumber`, gets the sequence number of the latest checkpoint that has been executed.
func (s *suiReadSystemFromSuiImpl) SuiGetLatestCheckpointSequenceNumber(ctx context.Context) (uint64, error) {
	var rsp uint64
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "sui_getLatestCheckpointSequenceNumber",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiXGetReferenceGasPrice implements the method `suix_getReferenceGasPrice`, gets the reference gas price for the network.
func (s *suiReadSystemFromSuiImpl) SuiXGetReferenceGasPrice(ctx context.Context) (uint64, error) {
	var rsp uint64
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getReferenceGasPrice",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiXGetCommitteeInfo implements the method `suix_getCommitteeInfo`, gets the committee information for the asked `epoch`.
func (s *suiReadSystemFromSuiImpl) SuiXGetCommitteeInfo(ctx context.Context, req models.SuiXGetCommitteeInfoRequest) (models.SuiXGetCommitteeInfoResponse, error) {
	var rsp models.SuiXGetCommitteeInfoResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getCommitteeInfo",
		Params: []interface{}{
			req.Epoch,
		},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiXGetStakes implements the method `suix_getStakes`, gets the delegated stakes for an address.
func (s *suiReadSystemFromSuiImpl) SuiXGetStakes(ctx context.Context, req models.SuiXGetStakesRequest) ([]*models.DelegatedStakesResponse, error) {
	var rsp []*models.DelegatedStakesResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getStakes",
		Params: []interface{}{
			req.Owner,
		},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiXGetStakesByIds implements the method `suix_getStakesByIds`, return one or more delegated stake. If a Stake was withdrawn, its status will be Unstaked.
func (s *suiReadSystemFromSuiImpl) SuiXGetStakesByIds(ctx context.Context, req models.SuiXGetStakesByIdsRequest) ([]*models.DelegatedStakesResponse, error) {
	var rsp []*models.DelegatedStakesResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getStakesByIds",
		Params: []interface{}{
			req.StakedSuiIds,
		},
	})
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}

// SuiXGetEpochs implements the method `suix_getEpochs`, get a list of epoch info.
func (s *suiReadSystemFromSuiImpl) SuiXGetEpochs(ctx context.Context, req models.SuiXGetEpochsRequest) (models.PaginatedEpochInfoResponse, error) {
	var rsp models.PaginatedEpochInfoResponse
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getEpochs",
		Params: []interface{}{
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

// SuiXGetCurrentEpoch implements the method `suix_getCurrentEpoch`, get current epoch info.
func (s *suiReadSystemFromSuiImpl) SuiXGetCurrentEpoch(ctx context.Context) (models.EpochInfo, error) {
	var rsp models.EpochInfo
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getCurrentEpoch",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// SuiXGetLatestSuiSystemState implements the method `suix_getLatestSuiSystemState`, get the latest SUI system state object on-chain.
func (s *suiReadSystemFromSuiImpl) SuiXGetLatestSuiSystemState(ctx context.Context) (models.SuiSystemStateSummary, error) {
	var rsp models.SuiSystemStateSummary
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getLatestSuiSystemState",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// SuiGetChainIdentifier implements the method `sui_getChainIdentifier`, return the chain's identifier.
func (s *suiReadSystemFromSuiImpl) SuiGetChainIdentifier(ctx context.Context) (string, error) {
	var rsp string
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "sui_getChainIdentifier",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

// SuiXGetValidatorsApy implements the method `suix_getValidatorsApy`, return the validator APY.
func (s *suiReadSystemFromSuiImpl) SuiXGetValidatorsApy(ctx context.Context) (models.ValidatorsApy, error) {
	var rsp models.ValidatorsApy
	err := s.conn.CallContext(ctx, &rsp, httpconn.Operation{
		Method: "suix_getValidatorsApy",
		Params: []interface{}{},
	})
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

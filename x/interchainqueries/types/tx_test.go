package types_test

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"

	iqtypes "github.com/neutron-org/neutron/x/interchainqueries/types"
)

const TestAddress = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"

func TestMsgRegisterInterchainQueryValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"invalid query type",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          "invalid_type",
					UpdatePeriod:       1,
					Sender:             TestAddress,
				}
			},
			iqtypes.ErrInvalidQueryType,
		},
		{
			"invalid transactions filter format",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "&)(^Y(*&(*&(&(*",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             TestAddress,
				}
			},
			iqtypes.ErrInvalidTransactionsFilter,
		},
		{
			"invalid update period",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       0,
					Sender:             TestAddress,
				}
			},
			iqtypes.ErrInvalidUpdatePeriod,
		},
		{
			"empty sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             "",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             "cosmos14234_invalid_address",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"empty connection id",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             TestAddress,
				}
			},
			iqtypes.ErrInvalidConnectionID,
		},
		{
			"valid",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod:       1,
					Sender:             TestAddress,
				}
			},
			nil,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgSubmitQueryResultValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"valid",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   TestAddress,
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key: []byte{10},
							Proof: &crypto.ProofOps{Ops: []crypto.ProofOp{
								{
									Type: "type",
									Key:  []byte{10},
									Data: []byte{10},
								},
							}},
							Value:         []byte{10},
							StoragePrefix: host.StoreKey,
						}},
						Block:    nil,
						Height:   100,
						Revision: 1,
					},
				}
			},
			nil,
		},
		{
			"empty result",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   TestAddress,
					ClientId: "client-id",
					Result:   nil,
				}
			},
			iqtypes.ErrEmptyResult,
		},
		{
			"empty kv results and block result",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   TestAddress,
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: nil,
						Block:     nil,
						Height:    100,
						Revision:  1,
					},
				}
			},
			iqtypes.ErrEmptyResult,
		},
		{
			"zero query id",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  0,
					Sender:   TestAddress,
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key: []byte{10},
							Proof: &crypto.ProofOps{Ops: []crypto.ProofOp{
								{
									Type: "type",
									Key:  []byte{10},
									Data: []byte{10},
								},
							}},
							Value:         []byte{10},
							StoragePrefix: host.StoreKey,
						}},
						Block:    nil,
						Height:   100,
						Revision: 1,
					},
				}
			},
			iqtypes.ErrInvalidQueryID,
		},
		{
			"empty sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   "",
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key: []byte{10},
							Proof: &crypto.ProofOps{Ops: []crypto.ProofOp{
								{
									Type: "type",
									Key:  []byte{10},
									Data: []byte{10},
								},
							}},
							Value:         []byte{10},
							StoragePrefix: host.StoreKey,
						}},
						Block:    nil,
						Height:   100,
						Revision: 1,
					},
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   "invalid_sender",
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key: []byte{10},
							Proof: &crypto.ProofOps{Ops: []crypto.ProofOp{
								{
									Type: "type",
									Key:  []byte{10},
									Data: []byte{10},
								},
							}},
							Value:         []byte{10},
							StoragePrefix: host.StoreKey,
						}},
						Block:    nil,
						Height:   100,
						Revision: 1,
					},
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"empty client id",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   TestAddress,
					ClientId: "",
					Result: &iqtypes.QueryResult{
						KvResults: nil,
						Block: &iqtypes.Block{
							NextBlockHeader: nil,
							Header:          nil,
							Tx:              nil,
						},
						Height:   100,
						Revision: 1,
					},
				}
			},
			iqtypes.ErrInvalidClientID,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgUpdateQueryRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"valid",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 1,
					NewKeys: []*iqtypes.KVKey{{
						Path: "staking",
						Key:  []byte{1, 2, 3},
					}},
					NewUpdatePeriod: 10,
					Sender:          TestAddress,
				}
			},
			nil,
		},
		{
			"empty keys and update_period",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId:         1,
					NewKeys:         nil,
					NewUpdatePeriod: 0,
					Sender:          TestAddress,
				}
			},
			sdkerrors.ErrInvalidRequest,
		},
		{
			"invalid query id",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 0,
					NewKeys: []*iqtypes.KVKey{{
						Path: "staking",
						Key:  []byte{1, 2, 3},
					}},
					NewUpdatePeriod: 10,
					Sender:          TestAddress,
				}
			},
			iqtypes.ErrInvalidQueryID,
		},
		{
			"empty sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 1,
					NewKeys: []*iqtypes.KVKey{{
						Path: "staking",
						Key:  []byte{1, 2, 3},
					}},
					NewUpdatePeriod: 10,
					Sender:          "",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 1,
					NewKeys: []*iqtypes.KVKey{{
						Path: "staking",
						Key:  []byte{1, 2, 3},
					}},
					NewUpdatePeriod: 10,
					Sender:          "invalid-sender",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgRemoveQueryRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"valid",
			func() sdktypes.Msg {
				return &iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 1,
					Sender:  TestAddress,
				}
			},
			nil,
		},
		{
			"invalid query id",
			func() sdktypes.Msg {
				return &iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 0,
					Sender:  TestAddress,
				}
			},
			iqtypes.ErrInvalidQueryID,
		},
		{
			"empty sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 1,
					Sender:  "",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid sender",
			func() sdktypes.Msg {
				return &iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 1,
					Sender:  "invalid-sender",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgRegisterInterchainQueryGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "connection-0",
					TransactionsFilter: "{}",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             TestAddress,
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}

func TestMsgSubmitQueryResultGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   TestAddress,
					ClientId: "client-id",
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key: []byte{10},
							Proof: &crypto.ProofOps{Ops: []crypto.ProofOp{
								{
									Type: "type",
									Key:  []byte{10},
									Data: []byte{10},
								},
							}},
							Value:         []byte{10},
							StoragePrefix: host.StoreKey,
						}},
						Block:    nil,
						Height:   100,
						Revision: 1,
					},
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}

func TestMsgUpdateQueryGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &iqtypes.MsgUpdateInterchainQueryRequest{
					Sender: TestAddress,
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}

func TestMsgRemoveQueryGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &iqtypes.MsgRemoveInterchainQueryRequest{
					Sender: TestAddress,
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}

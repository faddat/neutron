package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/neutron-org/neutron/testutil"
	"github.com/neutron-org/neutron/x/interchainqueries/keeper"
	iqtypes "github.com/neutron-org/neutron/x/interchainqueries/types"
)

var reflectContractPath = "../../../wasmbinding/testdata/reflect.wasm"

type KeeperTestSuite struct {
	testutil.IBCConnectionTestSuite
}

func (suite *KeeperTestSuite) TestRegisterInterchainQuery() {
	var msg iqtypes.MsgRegisterInterchainQuery

	tests := []struct {
		name         string
		topupBalance bool
		malleate     func(sender string)
		expectedErr  error
	}{
		{
			"invalid connection",
			true,
			func(sender string) {
				msg = iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       "unknown",
					TransactionsFilter: "[]",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             sender,
				}
			},
			iqtypes.ErrInvalidConnectionID,
		},
		{
			"insufficient funds for deposit",
			false,
			func(sender string) {
				msg = iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       suite.Path.EndpointA.ConnectionID,
					TransactionsFilter: "[]",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             sender,
				}
			},
			sdkerrors.ErrInsufficientFunds,
		},
		{
			"valid",
			true,
			func(sender string) {
				msg = iqtypes.MsgRegisterInterchainQuery{
					ConnectionId:       suite.Path.EndpointA.ConnectionID,
					TransactionsFilter: "[]",
					Keys:               nil,
					QueryType:          string(iqtypes.InterchainQueryTypeTX),
					UpdatePeriod:       1,
					Sender:             sender,
				}
			},
			nil,
		},
	}

	for _, tt := range tests {
		suite.SetupTest()

		var (
			ctx           = suite.ChainA.GetContext()
			contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
		)

		// Store code and instantiate reflect contract.
		codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
		contractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
		suite.Require().NotEmpty(contractAddress)

		err := testutil.SetupICAPath(suite.Path, contractAddress.String())
		suite.Require().NoError(err)

		tt.malleate(contractAddress.String())

		if tt.topupBalance {
			// Top up contract address with native coins for deposit
			senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
			suite.TopUpWallet(ctx, senderAddress, contractAddress)
		}

		msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

		res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &msg)

		if tt.expectedErr != nil {
			suite.Require().ErrorIs(err, tt.expectedErr)
			suite.Require().Nil(res)
		} else {
			query, _ := keeper.Keeper.RegisteredQuery(
				suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper, sdk.WrapSDKContext(ctx),
				&iqtypes.QueryRegisteredQueryRequest{QueryId: 1})

			suite.Require().Equal(iqtypes.DefaultQueryDeposit, query.RegisteredQuery.Deposit)
			suite.Require().Equal(iqtypes.DefaultQuerySubmitTimeout, query.RegisteredQuery.SubmitTimeout)
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		}
	}
}

func (suite *KeeperTestSuite) TestUpdateInterchainQuery() {
	var msg iqtypes.MsgUpdateInterchainQueryRequest
	originalQuery := iqtypes.MsgRegisterInterchainQuery{
		QueryType: string(iqtypes.InterchainQueryTypeKV),
		Keys: []*iqtypes.KVKey{
			{
				Path: "somepath",
				Key:  []byte("somedata"),
			},
		},
		TransactionsFilter: "",
		ConnectionId:       suite.Path.EndpointA.ConnectionID,
		UpdatePeriod:       1,
		Sender:             "",
	}

	tests := []struct {
		name              string
		malleate          func(sender string)
		expectedErr       error
		expectedPeriod    uint64
		expectedQueryKeys []*iqtypes.KVKey
	}{
		{
			"valid update period",
			func(sender string) {
				msg = iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId:         1,
					NewKeys:         nil,
					NewUpdatePeriod: 2,
					Sender:          sender,
				}
			},
			nil,
			2,
			originalQuery.Keys,
		},
		{
			"valid query data",
			func(sender string) {
				msg = iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 1,
					NewKeys: []*iqtypes.KVKey{
						{
							Path: "newpath",
							Key:  []byte("newdata"),
						},
					},
					NewUpdatePeriod: 0,
					Sender:          sender,
				}
			},
			nil,
			originalQuery.UpdatePeriod,
			[]*iqtypes.KVKey{
				{
					Path: "newpath",
					Key:  []byte("newdata"),
				},
			},
		},
		{
			"valid query both query keys and update period",
			func(sender string) {
				msg = iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 1,
					NewKeys: []*iqtypes.KVKey{
						{
							Path: "newpath",
							Key:  []byte("newdata"),
						},
					},
					NewUpdatePeriod: 2,
					Sender:          sender,
				}
			},
			nil,
			2,
			[]*iqtypes.KVKey{
				{
					Path: "newpath",
					Key:  []byte("newdata"),
				},
			},
		},
		{
			"invavid query id",
			func(sender string) {
				msg = iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId: 2,
					NewKeys: []*iqtypes.KVKey{
						{
							Path: "newpath",
							Key:  []byte("newdata"),
						},
					},
					NewUpdatePeriod: 2,
					Sender:          sender,
				}
			},
			iqtypes.ErrInvalidQueryID,
			originalQuery.UpdatePeriod,
			originalQuery.Keys,
		},
		{
			"failed due to auth error",
			func(sender string) {
				var (
					ctx           = suite.ChainA.GetContext()
					contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
				)
				codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
				newContractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
				suite.Require().NotEmpty(newContractAddress)
				msg = iqtypes.MsgUpdateInterchainQueryRequest{
					QueryId:         1,
					NewKeys:         nil,
					NewUpdatePeriod: 2,
					Sender:          newContractAddress.String(),
				}
			},
			sdkerrors.ErrUnauthorized,
			originalQuery.UpdatePeriod,
			originalQuery.Keys,
		},
	}

	for i, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tt.name, i, len(tests)), func() {
			suite.SetupTest()

			var (
				ctx           = suite.ChainA.GetContext()
				contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
			)

			// Store code and instantiate reflect contract.
			codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
			contractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
			suite.Require().NotEmpty(contractAddress)

			err := testutil.SetupICAPath(suite.Path, contractAddress.String())
			suite.Require().NoError(err)

			// Top up contract address with native coins for deposit
			senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
			suite.TopUpWallet(ctx, senderAddress, contractAddress)

			tt.malleate(contractAddress.String())

			iqkeeper := suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper

			msgSrv := keeper.NewMsgServerImpl(iqkeeper)

			originalQuery.Sender = contractAddress.String()
			resRegister, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &originalQuery)
			suite.Require().NoError(err)
			suite.Require().NotNil(resRegister)

			resUpdate, err := msgSrv.UpdateInterchainQuery(sdktypes.WrapSDKContext(ctx), &msg)

			if tt.expectedErr != nil {
				suite.Require().ErrorIs(err, tt.expectedErr)
				suite.Require().Nil(resUpdate)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resUpdate)
				updatedQuery, err := iqkeeper.GetQueryByID(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(tt.expectedQueryKeys, updatedQuery.GetKeys())
				suite.Require().Equal(tt.expectedPeriod, updatedQuery.GetUpdatePeriod())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRemoveInterchainQuery() {
	suite.SetupTest()

	var msg iqtypes.MsgRemoveInterchainQueryRequest
	originalQuery := iqtypes.MsgRegisterInterchainQuery{
		QueryType:          string(iqtypes.InterchainQueryTypeKV),
		Keys:               nil,
		TransactionsFilter: "",
		ConnectionId:       suite.Path.EndpointA.ConnectionID,
		UpdatePeriod:       1,
		Sender:             "",
	}

	tests := []struct {
		name        string
		malleate    func(sender string)
		expectedErr error
	}{
		{
			"valid remove",
			func(sender string) {
				msg = iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 1,
					Sender:  sender,
				}
			},
			nil,
		},
		{
			"invalid query id",
			func(sender string) {
				msg = iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 2,
					Sender:  sender,
				}
			},
			iqtypes.ErrInvalidQueryID,
		},
		{
			"failed due to auth error",
			func(sender string) {
				var (
					ctx           = suite.ChainA.GetContext()
					contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
				)
				codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
				newContractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
				suite.Require().NotEmpty(newContractAddress)
				msg = iqtypes.MsgRemoveInterchainQueryRequest{
					QueryId: 1,
					Sender:  newContractAddress.String(),
				}
			},
			sdkerrors.ErrUnauthorized,
		},
	}

	for i, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tt.name, i, len(tests)), func() {
			suite.SetupTest()

			var (
				ctx           = suite.ChainA.GetContext()
				contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
			)

			// Store code and instantiate reflect contract.
			codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
			contractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
			suite.Require().NotEmpty(contractAddress)

			err := testutil.SetupICAPath(suite.Path, contractAddress.String())
			suite.Require().NoError(err)

			// Top up contract address with native coins for deposit
			bankKeeper := suite.GetNeutronZoneApp(suite.ChainA).BankKeeper
			senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
			suite.TopUpWallet(ctx, senderAddress, contractAddress)

			tt.malleate(contractAddress.String())
			iqkeeper := suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper

			msgSrv := keeper.NewMsgServerImpl(iqkeeper)
			originalQuery.Sender = contractAddress.String()

			resRegister, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &originalQuery)
			suite.Require().NoError(err)
			suite.Require().NotNil(resRegister)

			balance, balanceErr := bankKeeper.Balance(
				sdktypes.WrapSDKContext(ctx),
				&banktypes.QueryBalanceRequest{
					Address: contractAddress.String(),
					Denom:   sdktypes.DefaultBondDenom,
				},
			)
			expectedCoin := sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(int64(0)))

			suite.Require().NoError(balanceErr)
			suite.Require().NotNil(balance)
			suite.Require().Equal(&expectedCoin, balance.Balance)

			clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)
			resp := suite.ChainB.App.Query(abci.RequestQuery{
				Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
				Height: suite.ChainB.LastHeader.Header.Height - 1,
				Data:   clientKey,
				Prove:  true,
			})

			err = iqkeeper.SaveKVQueryResult(ctx, 1, &iqtypes.QueryResult{
				KvResults: []*iqtypes.StorageValue{{
					Key:           resp.Key,
					Proof:         resp.ProofOps,
					Value:         resp.Value,
					StoragePrefix: host.StoreKey,
				}},
				Block:    nil,
				Height:   1,
				Revision: 1,
			})
			suite.Require().NoError(err)

			resUpdate, err := msgSrv.RemoveInterchainQuery(sdktypes.WrapSDKContext(ctx), &msg)

			if tt.expectedErr != nil {
				suite.Require().ErrorIs(err, tt.expectedErr)
				suite.Require().Nil(resUpdate)
				originalQuery, queryErr := iqkeeper.GetQueryByID(ctx, 1)
				suite.Require().NoError(queryErr)
				suite.Require().NotNil(originalQuery)

				qr, qrerr := iqkeeper.GetQueryResultByID(ctx, 1)
				suite.Require().NoError(qrerr)
				suite.Require().NotNil(qr)
			} else {
				balance, balanceErr := bankKeeper.Balance(
					sdktypes.WrapSDKContext(ctx),
					&banktypes.QueryBalanceRequest{
						Address: contractAddress.String(),
						Denom:   sdktypes.DefaultBondDenom,
					},
				)
				expectedCoin := sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(int64(1_000_000)))

				suite.Require().NoError(balanceErr)
				suite.Require().NotNil(balance)
				suite.Require().Equal(&expectedCoin, balance.Balance)

				suite.Require().NoError(err)
				suite.Require().NotNil(resUpdate)
				originalQuery, queryErr := iqkeeper.GetQueryByID(ctx, 1)
				suite.Require().Error(queryErr, iqtypes.ErrInvalidQueryID)
				suite.Require().Nil(originalQuery)

				qr, qrerr := iqkeeper.GetQueryResultByID(ctx, 1)
				suite.Require().Error(qrerr, iqtypes.ErrNoQueryResult)
				suite.Require().Nil(qr)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSubmitInterchainQueryResult() {
	var msg iqtypes.MsgSubmitQueryResult

	tests := []struct {
		name          string
		malleate      func(sender string, ctx sdktypes.Context)
		expectedError error
	}{
		{
			"invalid query id",
			func(sender string, ctx sdktypes.Context) {
				// now we don't care what is really under the value, we just need to be sure that we can verify KV proofs
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)
				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  1,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer, and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			iqtypes.ErrInvalidQueryID,
		},
		{
			"valid KV storage proof",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)
				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				// suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer,
						// and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			nil,
		},
		{
			"non-registered key in KV result",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)

				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				suite.NoError(suite.Path.EndpointA.UpdateClient())

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   []byte("non-registered key"),
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer, and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			iqtypes.ErrInvalidSubmittedResult,
		},
		{
			"non-registered path in KV result",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)

				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: "non-registered-path",
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer,
						// and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			iqtypes.ErrInvalidSubmittedResult,
		},
		{
			"non existence KV proof",
			func(sender string, ctx sdktypes.Context) {
				clientKey := []byte("non_existed_key")

				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				// suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				// now we don't care what is really under the value, we just need to be sure that we can verify KV proofs
				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender, // A bit weird that query owner submits the results, but it doesn't really matter
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer,
						// and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			nil,
		},
		{
			"header with invalid height",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)
				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer, and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			ibcclienttypes.ErrConsensusStateNotFound,
		},
		{
			"invalid KV storage value",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)
				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         []byte("some evil data"),
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer, and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			iqtypes.ErrInvalidProof,
		},
		{
			"query result height is too old",
			func(sender string, ctx sdktypes.Context) {
				clientKey := host.FullClientStateKey(suite.Path.EndpointB.ClientID)

				registerMsg := iqtypes.MsgRegisterInterchainQuery{
					ConnectionId: suite.Path.EndpointA.ConnectionID,
					Keys: []*iqtypes.KVKey{
						{Path: host.StoreKey, Key: clientKey},
					},
					QueryType:    string(iqtypes.InterchainQueryTypeKV),
					UpdatePeriod: 1,
					Sender:       sender,
				}

				msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

				res, err := msgSrv.RegisterInterchainQuery(sdktypes.WrapSDKContext(ctx), &registerMsg)
				suite.Require().NoError(err)

				suite.NoError(suite.Path.EndpointB.UpdateClient())
				suite.NoError(suite.Path.EndpointA.UpdateClient())

				// pretend like we have a very new query result
				suite.NoError(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper.UpdateLastRemoteHeight(ctx, res.Id, 9999))

				resp := suite.ChainB.App.Query(abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
					Height: suite.ChainB.LastHeader.Header.Height - 1,
					Data:   clientKey,
					Prove:  true,
				})

				msg = iqtypes.MsgSubmitQueryResult{
					QueryId:  res.Id,
					Sender:   sender,
					ClientId: suite.Path.EndpointA.ClientID,
					Result: &iqtypes.QueryResult{
						KvResults: []*iqtypes.StorageValue{{
							Key:           resp.Key,
							Proof:         resp.ProofOps,
							Value:         resp.Value,
							StoragePrefix: host.StoreKey,
						}},
						// we don't have tests to test transactions proofs verification since it's a tendermint layer, and we don't have access to it here
						Block:    nil,
						Height:   uint64(resp.Height),
						Revision: suite.ChainA.LastHeader.GetHeight().GetRevisionNumber(),
					},
				}
			},
			iqtypes.ErrInvalidHeight,
		},
	}

	for i, tc := range tests {
		tt := tc
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tt.name, i, len(tests)), func() {
			suite.SetupTest()

			var (
				ctx           = suite.ChainA.GetContext()
				contractOwner = wasmKeeper.RandomAccountAddress(suite.T())
			)

			// Store code and instantiate reflect contract.
			codeId := suite.StoreReflectCode(ctx, contractOwner, reflectContractPath)
			contractAddress := suite.InstantiateReflectContract(ctx, contractOwner, codeId)
			suite.Require().NotEmpty(contractAddress)

			err := testutil.SetupICAPath(suite.Path, contractAddress.String())
			suite.Require().NoError(err)

			// Top up contract address with native coins for deposit
			senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
			suite.TopUpWallet(ctx, senderAddress, contractAddress)

			tt.malleate(contractAddress.String(), ctx)

			msgSrv := keeper.NewMsgServerImpl(suite.GetNeutronZoneApp(suite.ChainA).InterchainQueriesKeeper)

			res, err := msgSrv.SubmitQueryResult(sdktypes.WrapSDKContext(ctx), &msg)

			if tt.expectedError != nil {
				suite.Require().ErrorIs(err, tt.expectedError)
				suite.Require().Nil(res)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TopUpWallet(ctx sdktypes.Context, sender sdktypes.AccAddress, contractAddress sdktypes.AccAddress) {
	coinsAmnt := sdktypes.NewCoins(sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(int64(1_000_000))))
	bankKeeper := suite.GetNeutronZoneApp(suite.ChainA).BankKeeper
	bankKeeper.SendCoins(ctx, sender, contractAddress, coinsAmnt)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

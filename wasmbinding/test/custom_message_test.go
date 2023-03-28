package test

import (
	"encoding/json"
	"testing"

	ictxtypes "github.com/neutron-org/neutron/x/interchaintxs/types"

	adminkeeper "github.com/cosmos/admin-module/x/adminmodule/keeper"
	admintypes "github.com/cosmos/admin-module/x/adminmodule/types"

	"github.com/neutron-org/neutron/app/params"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmvm/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/neutron-org/neutron/app"
	"github.com/neutron-org/neutron/testutil"
	"github.com/neutron-org/neutron/wasmbinding"
	"github.com/neutron-org/neutron/wasmbinding/bindings"
	feetypes "github.com/neutron-org/neutron/x/feerefunder/types"
	icqkeeper "github.com/neutron-org/neutron/x/interchainqueries/keeper"
	icqtypes "github.com/neutron-org/neutron/x/interchainqueries/types"
	ictxkeeper "github.com/neutron-org/neutron/x/interchaintxs/keeper"
)

type CustomMessengerTestSuite struct {
	testutil.IBCConnectionTestSuite
	neutron         *app.App
	ctx             sdk.Context
	messenger       *wasmbinding.CustomMessenger
	contractOwner   sdk.AccAddress
	contractAddress sdk.AccAddress
}

func (suite *CustomMessengerTestSuite) SetupTest() {
	suite.IBCConnectionTestSuite.SetupTest()
	suite.neutron = suite.GetNeutronZoneApp(suite.ChainA)
	suite.ctx = suite.ChainA.GetContext()
	suite.messenger = &wasmbinding.CustomMessenger{}
	suite.messenger.Ictxmsgserver = ictxkeeper.NewMsgServerImpl(suite.neutron.InterchainTxsKeeper)
	suite.messenger.Keeper = suite.neutron.InterchainTxsKeeper
	suite.messenger.Icqmsgserver = icqkeeper.NewMsgServerImpl(suite.neutron.InterchainQueriesKeeper)
	suite.messenger.Adminserver = adminkeeper.NewMsgServerImpl(suite.neutron.AdminmoduleKeeper)
	suite.contractOwner = keeper.RandomAccountAddress(suite.T())
}

func (suite *CustomMessengerTestSuite) TestRegisterInterchainAccount() {
	// Store code and instantiate reflect contract
	codeID := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeID)
	suite.Require().NotEmpty(suite.contractAddress)

	// Craft RegisterInterchainAccount message
	msg, err := json.Marshal(bindings.NeutronMsg{
		RegisterInterchainAccount: &bindings.RegisterInterchainAccount{
			ConnectionId:        suite.Path.EndpointA.ConnectionID,
			InterchainAccountId: testutil.TestInterchainID,
		},
	})
	suite.NoError(err)

	// Dispatch RegisterInterchainAccount message
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	suite.Equal([][]byte{[]byte(`{}`)}, data)
}

func (suite *CustomMessengerTestSuite) TestRegisterInterchainAccountLongID() {
	// Store code and instantiate reflect contract
	codeID := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeID)
	suite.Require().NotEmpty(suite.contractAddress)

	// Craft RegisterInterchainAccount message
	msg, err := json.Marshal(bindings.NeutronMsg{
		RegisterInterchainAccount: &bindings.RegisterInterchainAccount{
			ConnectionId: suite.Path.EndpointA.ConnectionID,
			// the limit is 47, this line is 50 characters long
			InterchainAccountId: "01234567890123456789012345678901234567890123456789",
		},
	})
	suite.NoError(err)

	// Dispatch RegisterInterchainAccount message
	_, _, err = suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.Error(err)
	suite.ErrorIs(err, ictxtypes.ErrLongInterchainAccountID)
}

func (suite *CustomMessengerTestSuite) TestRegisterInterchainQuery() {
	// Store code and instantiate reflect contract
	codeID := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeID)
	suite.Require().NotEmpty(suite.contractAddress)

	err := testutil.SetupICAPath(suite.Path, suite.contractAddress.String())
	suite.Require().NoError(err)

	// Top up contract balance
	senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
	coinsAmnt := sdk.NewCoins(sdk.NewCoin(params.DefaultDenom, sdk.NewInt(int64(10_000_000))))
	bankKeeper := suite.neutron.BankKeeper
	err = bankKeeper.SendCoins(suite.ctx, senderAddress, suite.contractAddress, coinsAmnt)
	suite.NoError(err)

	// Craft RegisterInterchainQuery message
	msg, err := json.Marshal(bindings.NeutronMsg{
		RegisterInterchainQuery: &bindings.RegisterInterchainQuery{
			QueryType: string(icqtypes.InterchainQueryTypeKV),
			Keys: []*icqtypes.KVKey{
				{Path: host.StoreKey, Key: host.FullClientStateKey(suite.Path.EndpointB.ClientID)},
			},
			TransactionsFilter: "{}",
			ConnectionId:       suite.Path.EndpointA.ConnectionID,
			UpdatePeriod:       20,
		},
	})
	suite.NoError(err)

	// Dispatch RegisterInterchainQuery message
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	suite.Equal([][]byte{[]byte(`{"id":1}`)}, data)
}

func (suite *CustomMessengerTestSuite) TestUpdateInterchainQuery() {
	// reuse register interchain query test to get query registered
	suite.TestRegisterInterchainQuery()

	// Craft UpdateInterchainQuery message
	msg, err := json.Marshal(bindings.NeutronMsg{
		UpdateInterchainQuery: &bindings.UpdateInterchainQuery{
			QueryId:         1,
			NewKeys:         nil,
			NewUpdatePeriod: 111,
		},
	})
	suite.NoError(err)

	// Dispatch UpdateInterchainQuery message
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	suite.Equal([][]byte{[]byte(`{}`)}, data)
}

func (suite *CustomMessengerTestSuite) TestUpdateInterchainQueryFailed() {
	// Craft UpdateInterchainQuery message
	msg, err := json.Marshal(bindings.NeutronMsg{
		UpdateInterchainQuery: &bindings.UpdateInterchainQuery{
			QueryId:         1,
			NewKeys:         nil,
			NewUpdatePeriod: 1,
		},
	})
	suite.NoError(err)

	// Dispatch UpdateInterchainQuery message
	owner, err := sdk.AccAddressFromBech32(testutil.TestOwnerAddress)
	suite.NoError(err)
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, owner, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	expectedErrMsg := "failed to update interchain query: failed to update interchain query: failed to get query by query id: there is no query with id: 1"
	suite.Require().ErrorContains(err, expectedErrMsg)
	suite.Nil(events)
	suite.Nil(data)
}

func (suite *CustomMessengerTestSuite) TestRemoveInterchainQuery() {
	// Reuse register interchain query test to get query registered
	suite.TestRegisterInterchainQuery()

	// Craft RemoveInterchainQuery message
	msg, err := json.Marshal(bindings.NeutronMsg{
		RemoveInterchainQuery: &bindings.RemoveInterchainQuery{
			QueryId: 1,
		},
	})
	suite.NoError(err)

	// Dispatch RemoveInterchainQuery message
	suite.NoError(err)
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	suite.Equal([][]byte{[]byte(`{}`)}, data)
}

func (suite *CustomMessengerTestSuite) TestRemoveInterchainQueryFailed() {
	// Craft RemoveInterchainQuery message
	msg, err := json.Marshal(bindings.NeutronMsg{
		RemoveInterchainQuery: &bindings.RemoveInterchainQuery{
			QueryId: 1,
		},
	})
	suite.NoError(err)

	// Dispatch RemoveInterchainQuery message
	owner, err := sdk.AccAddressFromBech32(testutil.TestOwnerAddress)
	suite.NoError(err)
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, owner, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	expectedErrMsg := "failed to remove interchain query: failed to remove interchain query: failed to get query by query id: there is no query with id: 1"
	suite.Require().ErrorContains(err, expectedErrMsg)
	suite.Nil(events)
	suite.Nil(data)
}

func (suite *CustomMessengerTestSuite) TestSubmitTx() {
	// Store code and instantiate reflect contract
	codeID := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeID)
	suite.Require().NotEmpty(suite.contractAddress)

	senderAddress := suite.ChainA.SenderAccounts[0].SenderAccount.GetAddress()
	coinsAmnt := sdk.NewCoins(sdk.NewCoin(params.DefaultDenom, sdk.NewInt(int64(10_000_000))))
	bankKeeper := suite.neutron.BankKeeper
	err := bankKeeper.SendCoins(suite.ctx, senderAddress, suite.contractAddress, coinsAmnt)
	suite.NoError(err)

	err = testutil.SetupICAPath(suite.Path, suite.contractAddress.String())
	suite.Require().NoError(err)

	events, data, err := suite.messenger.DispatchMsg(
		suite.ctx,
		suite.contractAddress,
		suite.Path.EndpointA.ChannelConfig.PortID,
		types.CosmosMsg{
			Custom: suite.craftMarshaledMsgSubmitTxWithNumMsgs(1),
		},
	)
	suite.NoError(err)

	var response bindings.SubmitTxResponse
	err = json.Unmarshal(data[0], &response)
	suite.NoError(err)
	suite.Nil(events)
	suite.Equal(uint64(1), response.SequenceId)
	suite.Equal("channel-2", response.Channel)
}

func (suite *CustomMessengerTestSuite) TestSubmitTxTooMuchTxs() {
	// Store code and instantiate reflect contract
	codeID := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeID)
	suite.Require().NotEmpty(suite.contractAddress)

	err := testutil.SetupICAPath(suite.Path, suite.contractAddress.String())
	suite.Require().NoError(err)

	_, _, err = suite.messenger.DispatchMsg(
		suite.ctx,
		suite.contractAddress,
		suite.Path.EndpointA.ChannelConfig.PortID,
		types.CosmosMsg{
			Custom: suite.craftMarshaledMsgSubmitTxWithNumMsgs(20),
		},
	)
	suite.ErrorContains(err, "MsgSubmitTx contains more messages than allowed")
}

func (suite *CustomMessengerTestSuite) TestSoftwareUpgradeProposal() {
	// Set admin so that we can execute this proposal without permission error
	suite.neutron.AdminmoduleKeeper.SetAdmin(suite.ctx, suite.contractAddress.String())

	// Craft SubmitAdminProposal message
	msg, err := json.Marshal(bindings.NeutronMsg{
		SubmitAdminProposal: &bindings.SubmitAdminProposal{
			AdminProposal: bindings.AdminProposal{
				SoftwareUpgradeProposal: &bindings.SoftwareUpgradeProposal{
					Title:       "Test",
					Description: "Test",
					Plan: bindings.Plan{
						Name:   "TestPlan",
						Height: 150,
						Info:   "TestInfo",
					},
				},
			},
		},
	})
	suite.NoError(err)

	// Dispatch SubmitAdminProposal message
	events, data, err := suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	expected, err := json.Marshal(&admintypes.MsgSubmitProposalResponse{
		ProposalId: 1,
	})
	suite.NoError(err)
	suite.Equal([][]uint8{expected}, data)

	// Test with other proposer that is not admin should return failure
	otherAddress, err := sdk.AccAddressFromBech32("neutron13jrwrtsyjjuynlug65r76r2zvfw5xjcq6532h2")
	suite.NoError(err)
	_, _, err = suite.messenger.DispatchMsg(suite.ctx, otherAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.Error(err)

	// Check CancelSubmitAdminProposal

	// Craft CancelSubmitAdminProposal message
	msg, err = json.Marshal(bindings.NeutronMsg{
		SubmitAdminProposal: &bindings.SubmitAdminProposal{
			AdminProposal: bindings.AdminProposal{
				CancelSoftwareUpgradeProposal: &bindings.CancelSoftwareUpgradeProposal{
					Title:       "Test",
					Description: "Test",
				},
			},
		},
	})
	suite.NoError(err)

	// Dispatch SubmitAdminProposal message
	events, data, err = suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, types.CosmosMsg{
		Custom: msg,
	})
	suite.NoError(err)
	suite.Nil(events)
	expected, err = json.Marshal(&admintypes.MsgSubmitProposalResponse{
		ProposalId: 2,
	})
	suite.NoError(err)
	suite.Equal([][]uint8{expected}, data)
}

func (suite *CustomMessengerTestSuite) TestTooMuchProposals() {
	// Store code and instantiate reflect contract
	codeId := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeId)
	suite.Require().NotEmpty(suite.contractAddress)

	err := testutil.SetupICAPath(suite.Path, suite.contractAddress.String())
	suite.Require().NoError(err)

	// Craft  message with 2 proposals
	msg, err := json.Marshal(bindings.NeutronMsg{
		SubmitAdminProposal: &bindings.SubmitAdminProposal{
			AdminProposal: bindings.AdminProposal{
				CancelSoftwareUpgradeProposal: &bindings.CancelSoftwareUpgradeProposal{
					Title:       "Test",
					Description: "Test",
				},
				ClearAdminProposal: &bindings.ClearAdminProposal{
					Title:       "Test",
					Description: "Test",
					Contract:    "Test",
				},
			},
		},
	})

	cosmosMsg := types.CosmosMsg{Custom: msg}

	// Dispatch SubmitAdminProposal message
	_, _, err = suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, cosmosMsg)

	suite.ErrorContains(err, "more than one admin proposal type is present in message")
}

func (suite *CustomMessengerTestSuite) TestNoProposals() {
	// Store code and instantiate reflect contract
	codeId := suite.StoreReflectCode(suite.ctx, suite.contractOwner, "../testdata/reflect.wasm")
	suite.contractAddress = suite.InstantiateReflectContract(suite.ctx, suite.contractOwner, codeId)
	suite.Require().NotEmpty(suite.contractAddress)

	err := testutil.SetupICAPath(suite.Path, suite.contractAddress.String())
	suite.Require().NoError(err)

	// Craft  message with 0 proposals
	msg, err := json.Marshal(bindings.NeutronMsg{
		SubmitAdminProposal: &bindings.SubmitAdminProposal{
			AdminProposal: bindings.AdminProposal{},
		},
	})

	cosmosMsg := types.CosmosMsg{Custom: msg}

	// Dispatch SubmitAdminProposal message
	_, _, err = suite.messenger.DispatchMsg(suite.ctx, suite.contractAddress, suite.Path.EndpointA.ChannelConfig.PortID, cosmosMsg)

	suite.ErrorContains(err, "no admin proposal type is present in message")
}

func (suite *CustomMessengerTestSuite) craftMarshaledMsgSubmitTxWithNumMsgs(numMsgs int) (result []byte) {
	msg := bindings.ProtobufAny{
		TypeURL: "/cosmos.staking.v1beta1.MsgDelegate",
		Value:   []byte{26, 10, 10, 5, 115, 116, 97, 107, 101, 18, 1, 48},
	}
	msgs := make([]bindings.ProtobufAny, 0, numMsgs)
	for i := 0; i < numMsgs; i++ {
		msgs = append(msgs, msg)
	}
	result, err := json.Marshal(struct {
		SubmitTx bindings.SubmitTx `json:"submit_tx"`
	}{
		SubmitTx: bindings.SubmitTx{
			ConnectionId:        suite.Path.EndpointA.ConnectionID,
			InterchainAccountId: testutil.TestInterchainID,
			Msgs:                msgs,
			Memo:                "Jimmy",
			Timeout:             2000,
			Fee: feetypes.Fee{
				RecvFee:    sdk.NewCoins(),
				AckFee:     sdk.NewCoins(sdk.NewCoin(params.DefaultDenom, sdk.NewInt(1000))),
				TimeoutFee: sdk.NewCoins(sdk.NewCoin(params.DefaultDenom, sdk.NewInt(1000))),
			},
		},
	})
	suite.NoError(err)
	return
}

func TestMessengerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomMessengerTestSuite))
}

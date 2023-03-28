package wasmbinding

import (
	"encoding/json"
	"fmt"

	paramChange "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	softwareUpgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	adminkeeper "github.com/cosmos/admin-module/x/adminmodule/keeper"
	admintypes "github.com/cosmos/admin-module/x/adminmodule/types"

	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	"github.com/neutron-org/neutron/wasmbinding/bindings"
	icqkeeper "github.com/neutron-org/neutron/x/interchainqueries/keeper"
	icqtypes "github.com/neutron-org/neutron/x/interchainqueries/types"
	ictxkeeper "github.com/neutron-org/neutron/x/interchaintxs/keeper"
	ictxtypes "github.com/neutron-org/neutron/x/interchaintxs/types"
	transferwrapperkeeper "github.com/neutron-org/neutron/x/transfer/keeper"
	transferwrappertypes "github.com/neutron-org/neutron/x/transfer/types"
)

func CustomMessageDecorator(ictx *ictxkeeper.Keeper, icq *icqkeeper.Keeper, transferKeeper transferwrapperkeeper.KeeperTransferWrapper, admKeeper *adminkeeper.Keeper) func(messenger wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			Keeper:         *ictx,
			Wrapped:        old,
			Ictxmsgserver:  ictxkeeper.NewMsgServerImpl(*ictx),
			Icqmsgserver:   icqkeeper.NewMsgServerImpl(*icq),
			transferKeeper: transferKeeper,
			Adminserver:    adminkeeper.NewMsgServerImpl(*admKeeper),
		}
	}
}

type CustomMessenger struct {
	Keeper         ictxkeeper.Keeper
	Wrapped        wasmkeeper.Messenger
	Ictxmsgserver  ictxtypes.MsgServer
	Icqmsgserver   icqtypes.MsgServer
	transferKeeper transferwrapperkeeper.KeeperTransferWrapper
	Adminserver    admintypes.MsgServer
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg bindings.NeutronMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			ctx.Logger().Debug("json.Unmarshal: failed to decode incoming custom cosmos message",
				"from_address", contractAddr.String(),
				"message", string(msg.Custom),
				"error", err,
			)
			return nil, nil, sdkerrors.Wrap(err, "failed to decode incoming custom cosmos message")
		}

		if contractMsg.SubmitTx != nil {
			return m.submitTx(ctx, contractAddr, contractMsg.SubmitTx)
		}
		if contractMsg.RegisterInterchainAccount != nil {
			return m.registerInterchainAccount(ctx, contractAddr, contractMsg.RegisterInterchainAccount)
		}
		if contractMsg.RegisterInterchainQuery != nil {
			return m.registerInterchainQuery(ctx, contractAddr, contractMsg.RegisterInterchainQuery)
		}
		if contractMsg.UpdateInterchainQuery != nil {
			return m.updateInterchainQuery(ctx, contractAddr, contractMsg.UpdateInterchainQuery)
		}
		if contractMsg.RemoveInterchainQuery != nil {
			return m.removeInterchainQuery(ctx, contractAddr, contractMsg.RemoveInterchainQuery)
		}
		if contractMsg.IBCTransfer != nil {
			return m.ibcTransfer(ctx, contractAddr, *contractMsg.IBCTransfer)
		}
		if contractMsg.SubmitAdminProposal != nil {
			return m.submitAdminProposal(ctx, contractAddr, contractMsg.SubmitAdminProposal)
		}
	}

	return m.Wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) ibcTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, ibcTransferMsg transferwrappertypes.MsgTransfer) ([]sdk.Event, [][]byte, error) {
	ibcTransferMsg.Sender = contractAddr.String()

	if err := ibcTransferMsg.ValidateBasic(); err != nil {
		return nil, nil, sdkerrors.Wrap(err, "failed to validate ibcTransferMsg")
	}

	response, err := m.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), &ibcTransferMsg)
	if err != nil {
		ctx.Logger().Debug("transferServer.Transfer: failed to transfer",
			"from_address", contractAddr.String(),
			"msg", ibcTransferMsg,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to execute IBCTransfer")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal MsgTransferResponse response to JSON",
			"from_address", contractAddr.String(),
			"msg", response,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("ibcTransferMsg completed",
		"from_address", contractAddr.String(),
		"msg", ibcTransferMsg,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) updateInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, updateQuery *bindings.UpdateInterchainQuery) ([]sdk.Event, [][]byte, error) {
	response, err := m.performUpdateInterchainQuery(ctx, contractAddr, updateQuery)
	if err != nil {
		ctx.Logger().Debug("performUpdateInterchainQuery: failed to update interchain query",
			"from_address", contractAddr.String(),
			"msg", updateQuery,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to update interchain query")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal UpdateInterchainQueryResponse response to JSON",
			"from_address", contractAddr.String(),
			"msg", updateQuery,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("interchain query updated",
		"from_address", contractAddr.String(),
		"msg", updateQuery,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) performUpdateInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, updateQuery *bindings.UpdateInterchainQuery) (*bindings.UpdateInterchainQueryResponse, error) {
	msg := icqtypes.MsgUpdateInterchainQueryRequest{
		QueryId:               updateQuery.QueryId,
		NewKeys:               updateQuery.NewKeys,
		NewUpdatePeriod:       updateQuery.NewUpdatePeriod,
		NewTransactionsFilter: updateQuery.NewTransactionsFilter,
		Sender:                contractAddr.String(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming UpdateInterchainQuery message")
	}

	response, err := m.Icqmsgserver.UpdateInterchainQuery(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update interchain query")
	}

	return (*bindings.UpdateInterchainQueryResponse)(response), nil
}

func (m *CustomMessenger) removeInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, removeQuery *bindings.RemoveInterchainQuery) ([]sdk.Event, [][]byte, error) {
	response, err := m.performRemoveInterchainQuery(ctx, contractAddr, removeQuery)
	if err != nil {
		ctx.Logger().Debug("performRemoveInterchainQuery: failed to update interchain query",
			"from_address", contractAddr.String(),
			"msg", removeQuery,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to remove interchain query")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal RemoveInterchainQueryResponse response to JSON",
			"from_address", contractAddr.String(),
			"msg", removeQuery,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("interchain query removed",
		"from_address", contractAddr.String(),
		"msg", removeQuery,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) performRemoveInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, updateQuery *bindings.RemoveInterchainQuery) (*bindings.RemoveInterchainQueryResponse, error) {
	msg := icqtypes.MsgRemoveInterchainQueryRequest{
		QueryId: updateQuery.QueryId,
		Sender:  contractAddr.String(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming RemoveInterchainQuery message")
	}

	response, err := m.Icqmsgserver.RemoveInterchainQuery(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to remove interchain query")
	}

	return (*bindings.RemoveInterchainQueryResponse)(response), nil
}

func (m *CustomMessenger) submitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) ([]sdk.Event, [][]byte, error) {
	response, err := m.performSubmitTx(ctx, contractAddr, submitTx)
	if err != nil {
		ctx.Logger().Debug("performSubmitTx: failed to submit interchain transaction",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to submit interchain transaction")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal submitTx response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("interchain transaction submitted",
		"from_address", contractAddr.String(),
		"connection_id", submitTx.ConnectionId,
		"interchain_account_id", submitTx.InterchainAccountId,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) submitAdminProposal(ctx sdk.Context, contractAddr sdk.AccAddress, submitAdminProposal *bindings.SubmitAdminProposal) ([]sdk.Event, [][]byte, error) {
	response, err := m.performSubmitAdminProposal(ctx, contractAddr, submitAdminProposal)
	if err != nil {
		ctx.Logger().Debug("performSubmitAdminProposal: failed to submitAdminProposal",
			"from_address", contractAddr.String(),
			"creator", contractAddr.String(),
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to submit admin proposal")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal submitAdminProposal response to JSON",
			"from_address", contractAddr.String(),
			"creator", contractAddr.String(),
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("submit proposal message submitted",
		"from_address", contractAddr.String(),
		"creator", contractAddr.String(),
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) performSubmitAdminProposal(ctx sdk.Context, contractAddr sdk.AccAddress, submitAdminProposal *bindings.SubmitAdminProposal) (*admintypes.MsgSubmitProposalResponse, error) {
	msg := admintypes.MsgSubmitProposal{Proposer: contractAddr.String()}
	proposal := submitAdminProposal.AdminProposal

	err := m.validateProposalQty(&proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate proposal quantity")
	}
	if proposal.ParamChangeProposal != nil {
		p := proposal.ParamChangeProposal
		err := msg.SetContent(&paramChange.ParameterChangeProposal{
			Title:       p.Title,
			Description: p.Description,
			Changes:     p.ParamChanges,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on ParameterChangeProposal")
		}
	}

	if proposal.SoftwareUpgradeProposal != nil {
		p := proposal.SoftwareUpgradeProposal
		err := msg.SetContent(&softwareUpgrade.SoftwareUpgradeProposal{
			Title:       p.Title,
			Description: p.Description,
			Plan: softwareUpgrade.Plan{
				Name:   p.Plan.Name,
				Height: p.Plan.Height,
				Info:   p.Plan.Info,
			},
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on SoftwareUpgradeProposal")
		}
	}

	if proposal.CancelSoftwareUpgradeProposal != nil {
		p := proposal.CancelSoftwareUpgradeProposal
		err := msg.SetContent(&softwareUpgrade.CancelSoftwareUpgradeProposal{
			Title:       p.Title,
			Description: p.Description,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on CancelSoftwareUpgradeProposal")
		}
	}

	if proposal.UpgradeProposal != nil {
		p := proposal.UpgradeProposal
		err := msg.SetContent(&ibcclienttypes.UpgradeProposal{
			Title:       p.Title,
			Description: p.Description,
			Plan: softwareUpgrade.Plan{
				Name:   p.Plan.Name,
				Height: p.Plan.Height,
				Info:   p.Plan.Info,
			},
			UpgradedClientState: p.UpgradedClientState,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on UpgradeProposal")
		}
	}

	if proposal.ClientUpdateProposal != nil {
		p := proposal.ClientUpdateProposal
		err := msg.SetContent(&ibcclienttypes.ClientUpdateProposal{
			Title:              p.Title,
			Description:        p.Description,
			SubjectClientId:    p.SubjectClientId,
			SubstituteClientId: p.SubstituteClientId,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on ClientUpdateProposal")
		}
	}

	if proposal.PinCodesProposal != nil {
		p := proposal.PinCodesProposal
		err := msg.SetContent(&wasmtypes.PinCodesProposal{
			Title:       p.Title,
			Description: p.Description,
			CodeIDs:     p.CodeIDs,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on PinCodesProposal")
		}
	}

	if proposal.UnpinCodesProposal != nil {
		p := proposal.UnpinCodesProposal
		err := msg.SetContent(&wasmtypes.UnpinCodesProposal{
			Title:       p.Title,
			Description: p.Description,
			CodeIDs:     p.CodeIDs,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on UnpinCodesProposal")
		}
	}

	if proposal.UpdateAdminProposal != nil {
		p := proposal.UpdateAdminProposal
		err := msg.SetContent(&wasmtypes.UpdateAdminProposal{
			Title:       p.Title,
			Description: p.Description,
			NewAdmin:    p.NewAdmin,
			Contract:    p.Contract,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on UpdateAdminProposal")
		}
	}

	if proposal.ClearAdminProposal != nil {
		p := proposal.ClearAdminProposal
		err := msg.SetContent(&wasmtypes.ClearAdminProposal{
			Title:       p.Title,
			Description: p.Description,
			Contract:    p.Contract,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to set content on ClearAdminProposal")
		}
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming SubmitAdminProposal message")
	}

	response, err := m.Adminserver.SubmitProposal(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to submit proposal")
	}

	return response, nil
}

func (m *CustomMessenger) performSubmitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) (*bindings.SubmitTxResponse, error) {
	tx := ictxtypes.MsgSubmitTx{
		FromAddress:         contractAddr.String(),
		ConnectionId:        submitTx.ConnectionId,
		Memo:                submitTx.Memo,
		InterchainAccountId: submitTx.InterchainAccountId,
		Timeout:             submitTx.Timeout,
		Fee:                 submitTx.Fee,
	}
	for _, msg := range submitTx.Msgs {
		tx.Msgs = append(tx.Msgs, &types.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		})
	}

	if err := tx.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming SubmitTx message")
	}

	response, err := m.Ictxmsgserver.SubmitTx(sdk.WrapSDKContext(ctx), &tx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to submit interchain transaction")
	}

	return (*bindings.SubmitTxResponse)(response), nil
}

func (m *CustomMessenger) registerInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) ([]sdk.Event, [][]byte, error) {
	response, err := m.performRegisterInterchainAccount(ctx, contractAddr, reg)
	if err != nil {
		ctx.Logger().Debug("performRegisterInterchainAccount: failed to register interchain account",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to register interchain account")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal register interchain account response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("registered interchain account",
		"from_address", contractAddr.String(),
		"connection_id", reg.ConnectionId,
		"interchain_account_id", reg.InterchainAccountId,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) performRegisterInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) (*bindings.RegisterInterchainAccountResponse, error) {
	msg := ictxtypes.MsgRegisterInterchainAccount{
		FromAddress:         contractAddr.String(),
		ConnectionId:        reg.ConnectionId,
		InterchainAccountId: reg.InterchainAccountId,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming RegisterInterchainAccount message")
	}

	response, err := m.Ictxmsgserver.RegisterInterchainAccount(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register interchain account")
	}

	return (*bindings.RegisterInterchainAccountResponse)(response), nil
}

func (m *CustomMessenger) registerInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainQuery) ([]sdk.Event, [][]byte, error) {
	response, err := m.performRegisterInterchainQuery(ctx, contractAddr, reg)
	if err != nil {
		ctx.Logger().Debug("performRegisterInterchainQuery: failed to register interchain query",
			"from_address", contractAddr.String(),
			"query_type", reg.QueryType,
			"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
			"transactions_filter", reg.TransactionsFilter,
			"connection_id", reg.ConnectionId,
			"update_period", reg.UpdatePeriod,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to register interchain query")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Error("json.Marshal: failed to marshal register interchain query response to JSON",
			"from_address", contractAddr.String(),
			"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
			"transactions_filter", reg.TransactionsFilter,
			"connection_id", reg.ConnectionId,
			"update_period", reg.UpdatePeriod,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("registered interchain query",
		"from_address", contractAddr.String(),
		"query_type", reg.QueryType,
		"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
		"transactions_filter", reg.TransactionsFilter,
		"connection_id", reg.ConnectionId,
		"update_period", reg.UpdatePeriod,
		"query_id", response.Id,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) performRegisterInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainQuery) (*bindings.RegisterInterchainQueryResponse, error) {
	msg := icqtypes.MsgRegisterInterchainQuery{
		Keys:               reg.Keys,
		TransactionsFilter: reg.TransactionsFilter,
		QueryType:          reg.QueryType,
		ConnectionId:       reg.ConnectionId,
		UpdatePeriod:       reg.UpdatePeriod,
		Sender:             contractAddr.String(),
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming RegisterInterchainQuery message")
	}

	response, err := m.Icqmsgserver.RegisterInterchainQuery(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register interchain query")
	}

	return (*bindings.RegisterInterchainQueryResponse)(response), nil
}

func (m *CustomMessenger) validateProposalQty(proposal *bindings.AdminProposal) error {
	qty := 0
	if proposal.ParamChangeProposal != nil {
		qty++
	}
	if proposal.SoftwareUpgradeProposal != nil {
		qty++
	}
	if proposal.CancelSoftwareUpgradeProposal != nil {
		qty++
	}
	if proposal.ClientUpdateProposal != nil {
		qty++
	}
	if proposal.UpgradeProposal != nil {
		qty++
	}
	if proposal.PinCodesProposal != nil {
		qty++
	}
	if proposal.UnpinCodesProposal != nil {
		qty++
	}
	if proposal.UpdateAdminProposal != nil {
		qty++
	}
	if proposal.ClearAdminProposal != nil {
		qty++
	}

	if qty == 0 {
		return fmt.Errorf("no admin proposal type is present in message")
	}

	if qty == 1 {
		return nil
	}

	return fmt.Errorf("more than one admin proposal type is present in message")
}

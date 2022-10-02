package types

import (
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
)

var _ codectypes.UnpackInterfacesMessage = MsgSubmitQueryResult{}

func (msg MsgSubmitQueryResult) Route() string {
	return RouterKey
}

func (msg MsgSubmitQueryResult) Type() string {
	return "submit-query-result"
}

func (msg MsgSubmitQueryResult) ValidateBasic() error {
	if msg.Result == nil {
		return sdkerrors.Wrap(ErrEmptyResult, "query result can't be empty")
	}

	if len(msg.Result.KvResults) == 0 && msg.Result.Block == nil {
		return sdkerrors.Wrap(ErrEmptyResult, "query result can't be empty")
	}

	if msg.QueryId == 0 {
		return sdkerrors.Wrap(ErrInvalidQueryID, "query id cannot be equal zero")
	}

	if strings.TrimSpace(msg.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Sender)
	}

	if strings.TrimSpace(msg.ClientId) == "" {
		return sdkerrors.Wrap(ErrInvalidClientID, "client id cannot be empty")
	}

	return nil
}

func (msg MsgSubmitQueryResult) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgSubmitQueryResult) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

func (msg MsgRegisterInterchainQuery) Route() string {
	return RouterKey
}

func (msg MsgRegisterInterchainQuery) Type() string {
	return "register-interchain-query"
}

func (msg MsgRegisterInterchainQuery) ValidateBasic() error {
	if msg.UpdatePeriod == 0 {
		return sdkerrors.Wrap(ErrInvalidUpdatePeriod, "update period can not be equal to zero")
	}

	if strings.TrimSpace(msg.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Sender)
	}

	if strings.TrimSpace(msg.ConnectionId) == "" {
		return sdkerrors.Wrap(ErrInvalidConnectionID, "connection id cannot be empty")
	}

	if !InterchainQueryType(msg.QueryType).IsValid() {
		return sdkerrors.Wrap(ErrInvalidQueryType, "invalid query type")
	}

	if InterchainQueryType(msg.QueryType).IsTX() {
		if err := ValidateTransactionsFilter(msg.TransactionsFilter); err != nil {
			return sdkerrors.Wrap(ErrInvalidTransactionsFilter, err.Error())
		}
	}
	return nil
}

func (msg MsgRegisterInterchainQuery) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgRegisterInterchainQuery) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgSubmitQueryResult) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var header exported.Header
	if err := unpacker.UnpackAny(msg.Result.GetBlock().GetHeader(), &header); err != nil {
		return err
	}

	if err := unpacker.UnpackAny(msg.Result.GetBlock().GetNextBlockHeader(), &header); err != nil {
		return err
	}

	return nil
}

func (msg MsgRemoveInterchainQueryRequest) ValidateBasic() error {
	if msg.GetQueryId() == 0 {
		return sdkerrors.Wrap(ErrInvalidQueryID, "query_id cannot be empty or equal to 0")
	}

	if strings.TrimSpace(msg.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Sender)
	}

	return nil
}

func (msg MsgRemoveInterchainQueryRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgRemoveInterchainQueryRequest) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

func (msg MsgUpdateInterchainQueryRequest) ValidateBasic() error {
	if msg.GetQueryId() == 0 {
		return sdkerrors.Wrap(ErrInvalidQueryID, "query_id cannot be empty or equal to 0")
	}

	if len(msg.GetNewKeys()) == 0 && msg.GetNewUpdatePeriod() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "one of new_keys or new_update_period should be set")
	}

	if strings.TrimSpace(msg.Sender) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Sender)
	}
	return nil
}

func (msg MsgUpdateInterchainQueryRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateInterchainQueryRequest) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

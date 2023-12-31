package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// slashing message types
const (
	TypeMsgUnjail = "unjail"
)

// verify interface at compile time
var (
	_ sdk.Msg = &MsgUnjail{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgImpeach{}
)

// NewMsgUnjail creates a new MsgUnjail instance
//
//nolint:interfacer
func NewMsgUnjail(validatorAddr sdk.AccAddress) *MsgUnjail {
	return &MsgUnjail{
		ValidatorAddr: validatorAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUnjail) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUnjail) Type() string { return TypeMsgUnjail }

// GetSigners returns the expected signers for MsgUnjail.
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddr)
	return []sdk.AccAddress{valAddr}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does a sanity check on the provided message.
func (msg MsgUnjail) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddr); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("validator input address: %s", err)
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHexUnsafe(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := msg.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// NewMsgImpeach creates a new MsgImpeach instance
func NewMsgImpeach(valAddr, from sdk.AccAddress) *MsgImpeach {
	return &MsgImpeach{
		ValidatorAddress: valAddr.String(),
		From:             from.String(),
	}
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgImpeach) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromHexUnsafe(msg.From)
	return []sdk.AccAddress{fromAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgImpeach) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgImpeach) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid account address: %s", err)
	}

	if _, err := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	return nil
}

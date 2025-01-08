package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _, _, _, _ sdk.Msg = &MsgSetWA0GI{}, &MsgSetMintCap{}, &MsgMint{}, &MsgBurn{}

func (msg *MsgSetWA0GI) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgSetWA0GI) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}

	return nil
}

func (msg MsgSetWA0GI) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

func (msg *MsgSetMintCap) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgSetMintCap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}

	return nil
}

func (msg MsgSetMintCap) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgMint) ValidateBasic() error {
	// forbid mint from cosmos tx
	// it can only be called by wrapped a0gi from EVM
	return ErrTxForbidden
}

func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

func (msg *MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgBurn) ValidateBasic() error {
	// forbid burn from cosmos tx
	// it can only be called by wrapped a0gi from EVM
	return ErrTxForbidden
}

func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

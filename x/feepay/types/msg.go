package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgRegisterFeePayContract{}
	_ sdk.Msg = &MsgUnregisterFeePayContract{}
	_ sdk.Msg = &MsgFundFeePayContract{}
	_ sdk.Msg = &MsgUpdateFeePayContractWalletLimit{}
	_ sdk.Msg = &MsgUpdateParams{}
)

const (
	TypeMsgRegisterFeePayContract          = "register_feepay_contract"
	TypeMsgUnregisterFeePayContract        = "unregister_feepay_contract"
	TypeMsgFundFeePayContract              = "fund_feepay_contract"
	TypeMsgUpdateFeePayContractWalletLimit = "update_feepay_contract_wallet_limit"
	TypeMsgUpdateParams                    = "msg_update_params"
)

// Route returns the name of the module
func (msg MsgRegisterFeePayContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterFeePayContract) Type() string { return TypeMsgRegisterFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterFeePayContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.FeePayContract.ContractAddress); err != nil {
		return err
	}

	if msg.FeePayContract.WalletLimit > 1_000_000 {
		return ErrInvalidWalletLimit
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterFeePayContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterFeePayContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnregisterFeePayContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnregisterFeePayContract) Type() string { return TypeMsgUnregisterFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnregisterFeePayContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnregisterFeePayContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnregisterFeePayContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgFundFeePayContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgFundFeePayContract) Type() string { return TypeMsgFundFeePayContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgFundFeePayContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return err
	}

	if len(msg.Amount) != 1 {
		return ErrInvalidFuryaFundAmount
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgFundFeePayContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgFundFeePayContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUpdateFeePayContractWalletLimit) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateFeePayContractWalletLimit) Type() string {
	return TypeMsgUpdateFeePayContractWalletLimit
}

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateFeePayContractWalletLimit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return err
	}

	if msg.WalletLimit > 1_000_000 {
		return ErrInvalidWalletLimit
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateFeePayContractWalletLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateFeePayContractWalletLimit) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrTxForbidden = errorsmod.Register(ModuleName, 1, "cosmos tx forbidden")
)

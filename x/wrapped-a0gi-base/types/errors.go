package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrTxForbidden            = errorsmod.Register(ModuleName, 1, "cosmos tx forbidden")
	ErrInsufficientMintCap    = errorsmod.Register(ModuleName, 2, "insufficient mint cap")
	ErrInsufficientMintSupply = errorsmod.Register(ModuleName, 3, "insufficient mint supply")
)

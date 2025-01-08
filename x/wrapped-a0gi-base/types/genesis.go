package types

import "github.com/ethereum/go-ethereum/common"

const (
	DEFAULT_WRAPPED_A0GI = "0000000000000000000000000000000000000000"
)

// NewGenesisState returns a new genesis state object for the module.
func NewGenesisState(addr common.Address) *GenesisState {
	return &GenesisState{
		WrappedA0GiAddress: addr.Bytes(),
	}
}

// DefaultGenesisState returns the default genesis state for the module.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(common.HexToAddress(DEFAULT_WRAPPED_A0GI))
}

// Validate performs basic validation of genesis data.
func (gs GenesisState) Validate() error {
	return nil
}

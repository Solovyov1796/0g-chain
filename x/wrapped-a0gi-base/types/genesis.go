package types

// NewGenesisState returns a new genesis state object for the module.
func NewGenesisState(addr string) *GenesisState {
	return &GenesisState{
		WrappedA0GiAddress: addr,
	}
}

// DefaultGenesisState returns the default genesis state for the module.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState("0000000000000000000000000000000000000000")
}

// Validate performs basic validation of genesis data.
func (gs GenesisState) Validate() error {
	return ValidateHexAddress(gs.WrappedA0GiAddress)
}

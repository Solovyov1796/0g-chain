package cli

import (
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for the inflation module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the dasigners module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand()

	return cmd
}

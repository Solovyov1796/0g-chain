package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/0glabs/0g-chain/x/committee/client/cli"
)

// ProposalHandler is a struct containing handler funcs for submiting CommitteeChange/Delete proposal txs to the gov module through the cli or rest.
var ProposalHandler = govclient.NewProposalHandler(cli.GetGovCmdSubmitProposal)

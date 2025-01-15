package wrappeda0gibase_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/0glabs/0g-chain/app"
	wrappeda0gibase "github.com/0glabs/0g-chain/x/wrapped-a0gi-base"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/testutil"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
)

type GenesisTestSuite struct {
	testutil.Suite
}

func (suite *GenesisTestSuite) TestInitGenesis() {
	// Most genesis validation tests are located in the types directory. The 'invalid' test cases are
	// randomly selected subset of those tests.
	testCases := []struct {
		name       string
		genState   *types.GenesisState
		expectPass bool
	}{
		{
			name:       "default genesis",
			genState:   types.DefaultGenesisState(),
			expectPass: true,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Setup (note: suite.SetupTest is not run before every suite.Run)
			suite.App = app.NewTestApp()
			suite.Keeper = suite.App.GetWrappedA0GIBaseKeeper()
			suite.Ctx = suite.App.NewContext(true, tmproto.Header{})

			// Run
			var exportedGenState *types.GenesisState
			run := func() {
				wrappeda0gibase.InitGenesis(suite.Ctx, suite.Keeper, *tc.genState)
				exportedGenState = wrappeda0gibase.ExportGenesis(suite.Ctx, suite.Keeper)
			}
			if tc.expectPass {
				suite.Require().NotPanics(run)
			} else {
				suite.Require().Panics(run)
			}

			// Check
			if tc.expectPass {
				expectedJson, err := suite.App.AppCodec().MarshalJSON(tc.genState)
				suite.Require().NoError(err)
				actualJson, err := suite.App.AppCodec().MarshalJSON(exportedGenState)
				suite.Require().NoError(err)
				suite.Equal(expectedJson, actualJson)
			}
		})
	}
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

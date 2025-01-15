package keeper_test

import (
	"testing"

	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/testutil"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	testutil.Suite
}

func (s *KeeperTestSuite) TestSetWA0GIAddress() {
	testCases := []struct {
		name  string
		wa0gi common.Address
	}{
		{
			name:  "zero address",
			wa0gi: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			s.Keeper.SetWA0GIAddress(s.Ctx, tc.wa0gi)
			response, err := s.Keeper.GetWA0GI(s.Ctx, &types.GetWA0GIRequest{})
			s.Require().NoError(err)
			s.Require().Equal(common.BytesToAddress(response.Address), tc.wa0gi)
		})
	}
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

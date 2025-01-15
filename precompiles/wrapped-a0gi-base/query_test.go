package wrappeda0gibase_test

import (
	"math/big"

	wrappeda0gibaseprecompile "github.com/0glabs/0g-chain/precompiles/wrapped-a0gi-base"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (s *WrappedA0giBaseTestSuite) TestGetW0GI() {
	method := wrappeda0gibaseprecompile.WrappedA0GIBaseFunctionGetWA0GI

	testCases := []struct {
		name        string
		malleate    func() []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func() []byte {
				input, err := s.abi.Pack(
					method,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				wa0gi := out[0].(common.Address)
				s.Require().Equal(wa0gi, common.HexToAddress(types.DEFAULT_WRAPPED_A0GI))
				// fmt.Println(wa0gi)
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			bz, err := s.runTx(tc.malleate(), s.signerOne, 10000000)

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(bz)
				tc.postCheck(bz)
			}
		})
	}
}

func (s *WrappedA0giBaseTestSuite) TestMinterSupply() {
	method := wrappeda0gibaseprecompile.WrappedA0GIBaseFunctionMinterSupply
	govAccAddr := s.App.GetGovKeeper().GetGovernanceAccount(s.Ctx).GetAddress().String()

	testCases := []struct {
		name        string
		malleate    func() []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"non-empty",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerOne.Addr,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				wa0gi := out[0].(wrappeda0gibaseprecompile.Supply)
				s.Require().Equal(wa0gi.Cap, big.NewInt(8e18))
				s.Require().Equal(wa0gi.Total, big.NewInt(1e18))
				// fmt.Println(wa0gi)
			},
			100000,
			false,
			"",
		}, {
			"empty",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerTwo.Addr,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				supply := out[0].(wrappeda0gibaseprecompile.Supply)
				s.Require().Equal(supply.Cap.Bytes(), big.NewInt(0).Bytes())
				s.Require().Equal(supply.Total.Bytes(), big.NewInt(0).Bytes())
				// fmt.Println(wa0gi)
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			s.wa0gibasekeeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMinterCap{
				Authority: govAccAddr,
				Minter:    s.signerOne.Addr.Bytes(),
				Cap:       big.NewInt(8e18).Bytes(),
			})

			s.wa0gibasekeeper.Mint(sdk.WrapSDKContext(s.Ctx), &types.MsgMint{
				Minter: s.signerOne.Addr.Bytes(),
				Amount: big.NewInt(1e18).Bytes(),
			})

			bz, err := s.runTx(tc.malleate(), s.signerOne, 10000000)

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(bz)
				tc.postCheck(bz)
			}
		})
	}
}

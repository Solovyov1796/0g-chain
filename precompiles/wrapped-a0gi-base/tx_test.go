package wrappeda0gibase_test

import (
	"fmt"
	"math/big"

	wrappeda0gibaseprecompile "github.com/0glabs/0g-chain/precompiles/wrapped-a0gi-base"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *WrappedA0giBaseTestSuite) TestMint() {
	method := wrappeda0gibaseprecompile.WrappedA0GIBaseFunctionMint
	govAccAddr := s.App.GetGovKeeper().GetGovernanceAccount(s.Ctx).GetAddress().String()

	testCases := []struct {
		name        string
		malleate    func() []byte
		postCheck   func()
		gas         uint64
		expErr      bool
		errContains string
		isSignerOne bool
	}{
		{
			"success",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerOne.Addr,
					big.NewInt(1e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {
				supply, err := s.wa0gibasekeeper.MinterSupply(s.Ctx, &types.MinterSupplyRequest{
					Address: s.signerOne.Addr.Bytes(),
				})
				s.Assert().NoError(err)
				s.Require().Equal(supply.Cap, big.NewInt(8e18).Bytes())
				s.Require().Equal(supply.Supply, big.NewInt(1e18).Bytes())
				// fmt.Println(wa0gi)
			},
			100000,
			false,
			"",
			true,
		}, {
			"fail",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerOne.Addr,
					big.NewInt(9e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {},
			100000,
			true,
			"insufficient mint cap",
			true,
		}, {
			"invalid sender",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerTwo.Addr,
					big.NewInt(9e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {},
			100000,
			true,
			"sender is not WA0GI",
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			fmt.Println(s.signerOne.Addr)
			s.wa0gibasekeeper.SetWA0GIAddress(s.Ctx, s.signerOne.Addr)
			s.wa0gibasekeeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMinterCap{
				Authority: govAccAddr,
				Minter:    s.signerOne.Addr.Bytes(),
				Cap:       big.NewInt(8e18).Bytes(),
			})

			var err error
			if tc.isSignerOne {
				_, err = s.runTx(tc.malleate(), s.signerOne, 10000000)
			} else {
				_, err = s.runTx(tc.malleate(), s.signerTwo, 10000000)
			}

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}

func (s *WrappedA0giBaseTestSuite) TestBurn() {
	method := wrappeda0gibaseprecompile.WrappedA0GIBaseFunctionBurn
	govAccAddr := s.App.GetGovKeeper().GetGovernanceAccount(s.Ctx).GetAddress().String()

	testCases := []struct {
		name        string
		malleate    func() []byte
		postCheck   func()
		gas         uint64
		expErr      bool
		errContains string
		isSignerOne bool
	}{
		{
			"success",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerOne.Addr,
					big.NewInt(1e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {
				supply, err := s.wa0gibasekeeper.MinterSupply(s.Ctx, &types.MinterSupplyRequest{
					Address: s.signerOne.Addr.Bytes(),
				})
				s.Assert().NoError(err)
				s.Require().Equal(supply.Cap, big.NewInt(8e18).Bytes())
				s.Require().Equal(supply.Supply, big.NewInt(3e18).Bytes())
				// fmt.Println(wa0gi)
			},
			100000,
			false,
			"",
			true,
		}, {
			"fail",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerOne.Addr,
					big.NewInt(9e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {},
			100000,
			true,
			"insufficient mint supply",
			true,
		}, {
			"invalid sender",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					s.signerTwo.Addr,
					big.NewInt(9e18),
				)
				s.Assert().NoError(err)
				return input
			},
			func() {},
			100000,
			true,
			"sender is not WA0GI",
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			fmt.Println(s.signerOne.Addr)
			s.wa0gibasekeeper.SetWA0GIAddress(s.Ctx, s.signerOne.Addr)
			s.wa0gibasekeeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMinterCap{
				Authority: govAccAddr,
				Minter:    s.signerOne.Addr.Bytes(),
				Cap:       big.NewInt(8e18).Bytes(),
			})
			s.wa0gibasekeeper.Mint(sdk.WrapSDKContext(s.Ctx), &types.MsgMint{
				Minter: s.signerOne.Addr.Bytes(),
				Amount: big.NewInt(4e18).Bytes(),
			})

			var err error
			if tc.isSignerOne {
				_, err = s.runTx(tc.malleate(), s.signerOne, 10000000)
			} else {
				_, err = s.runTx(tc.malleate(), s.signerTwo, 10000000)
			}

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}

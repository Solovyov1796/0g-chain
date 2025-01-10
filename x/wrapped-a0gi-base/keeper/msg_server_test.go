package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	precisebanktypes "github.com/0glabs/0g-chain/x/precisebank/types"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/testutil"
	"github.com/0glabs/0g-chain/x/wrapped-a0gi-base/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type MsgServerTestSuite struct {
	testutil.Suite
}

func (s *MsgServerTestSuite) TestSetWA0GI() {
	govAccAddr := s.GovKeeper.GetGovernanceAccount(s.Ctx).GetAddress().String()
	testCases := []struct {
		name      string
		req       *types.MsgSetWA0GI
		expectErr bool
		errMsg    string
	}{
		{
			name: "invalid signer",
			req: &types.MsgSetWA0GI{
				Authority: s.Addresses[0].String(),
				Address:   common.HexToAddress("0x0000000000000000000000000000000000000001").Bytes(),
			},
			expectErr: true,
			errMsg:    "expected gov account as only signer for proposal message",
		},
		{
			name: "success",
			req: &types.MsgSetWA0GI{
				Authority: govAccAddr,
				Address:   common.HexToAddress("0x0000000000000000000000000000000000000001").Bytes(),
			},
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.Keeper.SetWA0GI(sdk.WrapSDKContext(s.Ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				response, err := s.Keeper.GetWA0GI(s.Ctx, &types.GetWA0GIRequest{})
				s.Require().NoError(err)
				s.Require().Equal(response.Address, tc.req.Address)
			}
		})
	}
}

func (s *MsgServerTestSuite) TestSetMinterCap() {
	testCases := []struct {
		name string
		caps []struct {
			account common.Address
			cap     *big.Int
		}
	}{
		{
			name: "success",
			caps: []struct {
				account common.Address
				cap     *big.Int
			}{
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000000"),
					cap:     big.NewInt(100000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					cap:     big.NewInt(200000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					cap:     big.NewInt(300000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000003"),
					cap:     big.NewInt(400000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					cap:     big.NewInt(500000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					cap:     big.NewInt(600000),
				},
				{
					account: common.HexToAddress("0x0000000000000000000000000000000000000000"),
					cap:     big.NewInt(700000),
				},
			},
		},
	}
	s.Run("invalid authority", func() {
		s.SetupTest()
		_, err := s.Keeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMintCap{
			Authority: s.Addresses[0].String(),
			Minter:    common.HexToAddress("0x0000000000000000000000000000000000000000").Bytes(),
			Cap:       big.NewInt(600000).Bytes(),
		})
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "expected gov account as only signer for proposal message")
	})
	govAccAddr := s.GovKeeper.GetGovernanceAccount(s.Ctx).GetAddress().String()
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			c := make(map[common.Address]*big.Int)
			for _, cap := range tc.caps {
				_, err := s.Keeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMintCap{
					Authority: govAccAddr,
					Minter:    cap.account.Bytes(),
					Cap:       cap.cap.Bytes(),
				})
				s.Require().NoError(err)
				response, err := s.Keeper.MinterSupply(s.Ctx, &types.MinterSupplyRequest{
					Address: cap.account.Bytes(),
				})
				s.Require().NoError(err)
				s.Require().Equal(new(big.Int).SetBytes(response.Cap), cap.cap)
				c[cap.account] = cap.cap
			}
			for account, cap := range c {
				response, err := s.Keeper.MinterSupply(s.Ctx, &types.MinterSupplyRequest{
					Address: account.Bytes(),
				})
				s.Require().NoError(err)
				s.Require().Equal(new(big.Int).SetBytes(response.Cap), cap)
			}
		})
	}
}

type MintBurn struct {
	IsMint  bool
	Minter  common.Address
	Amount  *big.Int
	Success bool
}

func (s *MsgServerTestSuite) TestSetMintBurn() {
	precisebankKeeper := s.App.GetPrecisebankKeeper()
	accountKeeper := s.App.GetAccountKeeper()
	moduleAcc := accountKeeper.GetModuleAccount(s.Ctx, types.ModuleName).GetAddress()
	govAccAddr := s.GovKeeper.GetGovernanceAccount(s.Ctx).GetAddress().String()

	minter1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	minter2 := common.HexToAddress("0x0000000000000000000000000000000000000002")

	// set mint cap of minter 1 to 8 a0gi
	_, err := s.Keeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMintCap{
		Authority: govAccAddr,
		Minter:    minter1.Bytes(),
		Cap:       big.NewInt(8e18).Bytes(),
	})
	s.Require().NoError(err)
	// set mint cap of minter 2 to 5 a0gi
	_, err = s.Keeper.SetMinterCap(sdk.WrapSDKContext(s.Ctx), &types.MsgSetMintCap{
		Authority: govAccAddr,
		Minter:    minter2.Bytes(),
		Cap:       big.NewInt(5e18).Bytes(),
	})
	s.Require().NoError(err)

	testCases := []MintBurn{
		// #0, failed burn
		{
			IsMint:  false,
			Minter:  minter1,
			Amount:  big.NewInt(1e18),
			Success: false,
		},
		// #1, mint 5 a0gi by minter 1
		{
			IsMint:  true,
			Minter:  minter1,
			Amount:  big.NewInt(5e18),
			Success: true,
		},
		// #2, burn 0.5 a0gi by minter 1
		{
			IsMint:  false,
			Minter:  minter1,
			Amount:  big.NewInt(5e17),
			Success: true,
		},
		// #3, mint 0.7 a0gi by minter 2
		{
			IsMint:  true,
			Minter:  minter2,
			Amount:  big.NewInt(7e17),
			Success: true,
		},
		// #4, mint 2 a0gi by minter 2
		{
			IsMint:  true,
			Minter:  minter2,
			Amount:  big.NewInt(2e18),
			Success: true,
		},
		// #5, burn 0.3 a0gi by minter 2
		{
			IsMint:  false,
			Minter:  minter1,
			Amount:  big.NewInt(3e17),
			Success: true,
		},
		// #6, failed to mint 4 a0gi by minter 1
		{
			IsMint:  true,
			Minter:  minter1,
			Amount:  big.NewInt(4e18),
			Success: false,
		},
		// #7, mint 3.5 a0gi by minter 1
		{
			IsMint:  true,
			Minter:  minter1,
			Amount:  big.NewInt(3e18 + 5e17),
			Success: true,
		},
	}
	minted := big.NewInt(0)
	supplied := make(map[common.Address]*big.Int)
	for id, c := range testCases {
		fmt.Println(id)
		if c.IsMint {
			_, err = s.Keeper.Mint(sdk.WrapSDKContext(s.Ctx), &types.MsgMint{
				Minter: c.Minter.Bytes(),
				Amount: c.Amount.Bytes(),
			})
		} else {
			_, err = s.Keeper.Burn(sdk.WrapSDKContext(s.Ctx), &types.MsgBurn{
				Minter: c.Minter.Bytes(),
				Amount: c.Amount.Bytes(),
			})
		}
		if c.Success {
			if c.IsMint {
				minted.Add(minted, c.Amount)
				if amt, ok := supplied[c.Minter]; ok {
					amt.Add(amt, c.Amount)
				} else {
					supplied[c.Minter] = new(big.Int).Set(c.Amount)
				}
			} else {
				minted.Sub(minted, c.Amount)
				if amt, ok := supplied[c.Minter]; ok {
					amt.Sub(amt, c.Amount)
				} else {
					supplied[c.Minter] = new(big.Int).Set(c.Amount)
				}
			}
			s.Require().NoError(err)
			response, err := s.Keeper.MinterSupply(s.Ctx, &types.MinterSupplyRequest{
				Address: c.Minter.Bytes(),
			})
			s.Require().NoError(err)
			s.Require().Equal(supplied[c.Minter].Bytes(), response.Supply)
		} else {
			s.Require().Error(err)
		}
		coins := precisebankKeeper.GetBalance(s.Ctx, moduleAcc, precisebanktypes.ExtendedCoinDenom)
		s.Require().Equal(coins.Amount.BigInt(), minted)
	}
}

func TestMsgServerSuite(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}

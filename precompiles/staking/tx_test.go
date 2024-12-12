package staking_test

import (
	"encoding/base64"
	"math/big"
	"time"

	"cosmossdk.io/math"
	stakingprecompile "github.com/0glabs/0g-chain/precompiles/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/statedb"
)

func (s *StakingTestSuite) TestCreateValidator() {
	method := stakingprecompile.StakingFunctionCreateValidator
	description := stakingprecompile.Description{
		Moniker:         "test node",
		Identity:        "test node identity",
		Website:         "http://test.node.com",
		SecurityContact: "test node security contract",
		Details:         "test node details",
	}
	commission := stakingprecompile.CommissionRates{
		Rate:          math.LegacyOneDec().BigInt(),
		MaxRate:       math.LegacyOneDec().BigInt(),
		MaxChangeRate: math.LegacyOneDec().BigInt(),
	}
	minSelfDelegation := big.NewInt(1)
	pubkey := "eh/aR8BGUBIYI/Ust0NVBxZafLDAm7344F9dKzZU+7g="
	value := big.NewInt(100000000)

	testCases := []struct {
		name          string
		malleate      func() []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(data []byte)
		expError      bool
		errContains   string
	}{
		{
			"fail - ErrPubKeyInvalidLength",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					description,
					commission,
					minSelfDelegation,
					s.signerOne.HexAddr,
					value,
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func([]byte) {},
			true,
			stakingprecompile.ErrPubKeyInvalidLength,
		},
		{
			"success",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					description,
					commission,
					minSelfDelegation,
					pubkey,
					value,
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(data []byte) {},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			bz, err := s.runTx(tc.malleate(), s.signerOne, 10000000)
			s.stakingKeeper.ApplyAndReturnValidatorSetUpdates(s.Ctx)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				// query the validator in the staking keeper
				validator := s.StakingKeeper.Validator(s.Ctx, s.signerOne.ValAddr)
				s.Require().NoError(err)

				s.Require().NotNil(validator, "expected validator not to be nil")
				tc.postCheck(bz)

				isBonded := validator.IsBonded()
				s.Require().Equal(true, isBonded, "expected validator bonded to be %t; got %t", true, isBonded)

				consPubKey, err := validator.ConsPubKey()
				s.Require().NoError(err)
				consPubKeyBase64 := base64.StdEncoding.EncodeToString(consPubKey.Bytes())
				s.Require().Equal(pubkey, consPubKeyBase64, "expected validator pubkey to be %s; got %s", pubkey, consPubKeyBase64)

				operator := validator.GetOperator()
				s.Require().Equal(s.signerOne.ValAddr, operator, "expected validator operator to be %s; got %s", s.signerOne.ValAddr, operator)

				commissionRate := validator.GetCommission()
				s.Require().Equal(commission.Rate.String(), commissionRate.BigInt().String(), "expected validator commission rate to be %s; got %s", commission.Rate.String(), commissionRate.String())

				valMinSelfDelegation := validator.GetMinSelfDelegation()
				s.Require().Equal(minSelfDelegation.String(), valMinSelfDelegation.String(), "expected validator min self delegation to be %s; got %s", minSelfDelegation.String(), valMinSelfDelegation.String())

				moniker := validator.GetMoniker()
				s.Require().Equal(description.Moniker, moniker, "expected validator moniker to be %s; got %s", description.Moniker, moniker)

				jailed := validator.IsJailed()
				s.Require().Equal(false, jailed, "expected validator jailed to be %t; got %t", false, jailed)
			}
		})
	}
}

func (s *StakingTestSuite) TestEditValidator() {
	method := stakingprecompile.StakingFunctionEditValidator
	description := stakingprecompile.Description{
		Moniker:         "test node",
		Identity:        "test node identity",
		Website:         "http://test.node.com",
		SecurityContact: "test node security contract",
		Details:         "test node details",
	}
	newRate := math.LegacyOneDec().BigInt()
	newRate.Div(newRate, big.NewInt(2))
	minSelfDelegation := big.NewInt(2)

	testCases := []struct {
		name          string
		malleate      func() []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(data []byte)
		expError      bool
		errContains   string
	}{
		{
			"success",
			func() []byte {
				input, err := s.abi.Pack(
					method,
					description,
					stakingprecompile.NullableUint{
						IsNull: false,
						Value:  newRate,
					},
					stakingprecompile.NullableUint{
						IsNull: true,
						Value:  math.LegacyOneDec().BigInt(),
					},
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(data []byte) {},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			s.setupValidator(s.signerOne)
			// move block time forward
			s.Ctx = s.Ctx.WithBlockTime(time.Now().Add(time.Hour * 100))
			s.Statedb = statedb.New(s.Ctx, s.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))

			bz, err := s.runTx(tc.malleate(), s.signerOne, 10000000)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				// query the validator in the staking keeper
				validator := s.StakingKeeper.Validator(s.Ctx, s.signerOne.ValAddr)
				s.Require().NoError(err)

				s.Require().NotNil(validator, "expected validator not to be nil")
				tc.postCheck(bz)

				isBonded := validator.IsBonded()
				s.Require().Equal(true, isBonded, "expected validator bonded to be %t; got %t", false, isBonded)

				operator := validator.GetOperator()
				s.Require().Equal(s.signerOne.ValAddr, operator, "expected validator operator to be %s; got %s", s.signerOne.ValAddr, operator)

				commissionRate := validator.GetCommission()
				s.Require().Equal(newRate.String(), commissionRate.BigInt().String(), "expected validator commission rate to be %s; got %s", newRate.String(), commissionRate.String())

				valMinSelfDelegation := validator.GetMinSelfDelegation()
				s.Require().Equal(big.NewInt(1).String(), valMinSelfDelegation.String(), "expected validator min self delegation to be %s; got %s", minSelfDelegation.String(), valMinSelfDelegation.String())

				moniker := validator.GetMoniker()
				s.Require().Equal(description.Moniker, moniker, "expected validator moniker to be %s; got %s", description.Moniker, moniker)

				jailed := validator.IsJailed()
				s.Require().Equal(false, jailed, "expected validator jailed to be %t; got %t", false, jailed)
			}
		})
	}
}

func (s *StakingTestSuite) TestDelegate() {
	method := stakingprecompile.StakingFunctionDelegate

	testCases := []struct {
		name          string
		malleate      func(valAddr common.Address) []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(valAddr sdk.ValAddress)
		expError      bool
		errContains   string
	}{
		{
			"success",
			func(valAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					valAddr,
					big.NewInt(1000000),
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(valAddr sdk.ValAddress) {
				d, found := s.stakingKeeper.GetDelegation(s.Ctx, s.signerOne.AccAddr, valAddr)
				s.Assert().EqualValues(found, true)
				s.Assert().EqualValues(d.ValidatorAddress, valAddr.String())
				s.Assert().EqualValues(d.DelegatorAddress, s.signerOne.AccAddr.String())

				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("delegation: %s\n", string(jsonData))
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(operatorAddress)), s.signerOne, 10000000)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				tc.postCheck(operatorAddress)
			}
		})
	}
}

func (s *StakingTestSuite) TestBeginRedelegate() {
	method := stakingprecompile.StakingFunctionBeginRedelegate

	testCases := []struct {
		name          string
		malleate      func(srcAddr, dstAddr common.Address) []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(data []byte, srcAddr, dstAddr sdk.ValAddress)
		expError      bool
		errContains   string
	}{
		{
			"success",
			func(srcAddr, dstAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					srcAddr,
					dstAddr,
					big.NewInt(1000000),
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(data []byte, srcAddr, dstAddr sdk.ValAddress) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")

				d, found := s.stakingKeeper.GetRedelegation(s.Ctx, s.signerOne.AccAddr, srcAddr, dstAddr)
				s.Assert().EqualValues(found, true)
				s.Assert().EqualValues(d.DelegatorAddress, s.signerOne.AccAddr.String())
				s.Assert().EqualValues(d.ValidatorSrcAddress, srcAddr.String())
				s.Assert().EqualValues(d.ValidatorDstAddress, dstAddr.String())

				completionTime := out[0].(*big.Int)
				params := s.stakingKeeper.GetParams(s.Ctx)
				s.Assert().EqualValues(completionTime.Int64(), s.Ctx.BlockHeader().Time.Add(params.UnbondingTime).UTC().Unix())
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("redelegation: %s\n", string(jsonData))
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)

			// move block time forward
			s.Ctx = s.Ctx.WithBlockTime(time.Now().Add(time.Hour * 100))
			s.Statedb = statedb.New(s.Ctx, s.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))

			s.setupValidator(s.signerOne)

			bz, err := s.runTx(tc.malleate(s.signerOne.Addr, common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz, s.signerOne.ValAddr, operatorAddress)
			}
		})
	}
}

func (s *StakingTestSuite) TestUndelegate() {
	method := stakingprecompile.StakingFunctionUndelegate

	testCases := []struct {
		name          string
		malleate      func(valAddr common.Address) []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(data []byte, valAddr sdk.ValAddress)
		expError      bool
		errContains   string
	}{
		{
			"success",
			func(valAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					valAddr,
					big.NewInt(1000000),
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(data []byte, valAddr sdk.ValAddress) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")

				d, found := s.stakingKeeper.GetUnbondingDelegation(s.Ctx, s.signerOne.AccAddr, valAddr)
				s.Assert().EqualValues(found, true)
				s.Assert().EqualValues(d.DelegatorAddress, s.signerOne.AccAddr.String())
				s.Assert().EqualValues(d.ValidatorAddress, valAddr.String())

				completionTime := out[0].(*big.Int)
				params := s.stakingKeeper.GetParams(s.Ctx)
				s.Assert().EqualValues(completionTime.Int64(), s.Ctx.BlockHeader().Time.Add(params.UnbondingTime).UTC().Unix())
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("redelegation: %s\n", string(jsonData))
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			// move block time forward
			s.Ctx = s.Ctx.WithBlockTime(time.Now().Add(time.Hour * 100))
			s.Statedb = statedb.New(s.Ctx, s.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))

			s.setupValidator(s.signerOne)

			bz, err := s.runTx(tc.malleate(s.signerOne.Addr), s.signerOne, 10000000)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz, s.signerOne.ValAddr)
			}
		})
	}
}

func (s *StakingTestSuite) TestCancelUnbondingDelegation() {
	method := stakingprecompile.StakingFunctionCancelUnbondingDelegation

	testCases := []struct {
		name          string
		malleate      func(valAddr common.Address, height *big.Int) []byte
		gas           uint64
		callerAddress *common.Address
		postCheck     func(valAddr sdk.ValAddress)
		expError      bool
		errContains   string
	}{
		{
			"success",
			func(valAddr common.Address, height *big.Int) []byte {
				input, err := s.abi.Pack(
					method,
					valAddr,
					big.NewInt(1),
					height,
				)
				s.Assert().NoError(err)
				return input
			},
			200000,
			nil,
			func(valAddr sdk.ValAddress) {
				_, found := s.stakingKeeper.GetUnbondingDelegation(s.Ctx, s.signerOne.AccAddr, valAddr)
				s.Assert().EqualValues(found, false)
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("redelegation: %s\n", string(jsonData))
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			// move block time forward
			s.Ctx = s.Ctx.WithBlockTime(time.Now().Add(time.Hour * 100))
			s.Statedb = statedb.New(s.Ctx, s.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))

			s.setupValidator(s.signerOne)
			// unbond
			_, err := s.stakingKeeper.Undelegate(s.Ctx, s.signerOne.AccAddr, s.signerOne.ValAddr, sdk.NewDec(1))
			s.Require().NoError(err)

			u, _ := s.stakingKeeper.GetUnbondingDelegation(s.Ctx, s.signerOne.AccAddr, s.signerOne.ValAddr)
			height := u.Entries[0].CreationHeight

			bz, err := s.runTx(tc.malleate(s.signerOne.Addr, big.NewInt(height)), s.signerOne, 10000000)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				tc.postCheck(s.signerOne.ValAddr)
			}
		})
	}
}

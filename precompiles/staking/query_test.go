package staking_test

import (
	"math/big"

	stakingprecompile "github.com/0glabs/0g-chain/precompiles/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
)

func (s *StakingTestSuite) TestValidators() {
	method := stakingprecompile.StakingFunctionValidators

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
					"",
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				validators := out[0].([]stakingprecompile.Validator)
				paginationResult := out[1].(stakingprecompile.PageResponse)
				s.Assert().EqualValues(3, len(validators))
				s.Assert().EqualValues(3, paginationResult.Total)
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			s.AddDelegation(s.signerOne.HexAddr, s.signerTwo.HexAddr, sdk.NewIntFromUint64(1000000))

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

func (s *StakingTestSuite) TestValidator() {
	method := stakingprecompile.StakingFunctionValidator

	testCases := []struct {
		name        string
		malleate    func(operatorAddress common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(operatorAddress common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					operatorAddress,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				operatorAddress, err := s.firstBondedValidator()
				s.Require().NoError(err)
				validator := out[0].(stakingprecompile.Validator)
				s.Require().EqualValues(validator.OperatorAddress, common.BytesToAddress(operatorAddress.Bytes()))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestValidatorDelegations() {
	method := stakingprecompile.StakingFunctionValidatorDelegations

	testCases := []struct {
		name        string
		malleate    func(operatorAddress common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(operatorAddress common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					operatorAddress,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				operatorAddress, err := s.firstBondedValidator()
				s.Require().NoError(err)
				delegations := out[0].([]stakingprecompile.DelegationResponse)
				d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
				s.Require().EqualValues(len(delegations), len(d))
				// jsonData, _ := json.MarshalIndent(delegations, "", "    ")
				// fmt.Printf("delegations: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestValidatorUnbondingDelegations() {
	method := stakingprecompile.StakingFunctionValidatorUnbondingDelegations

	testCases := []struct {
		name        string
		malleate    func(operatorAddress common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(operatorAddress common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					operatorAddress,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				unbonding := out[0].([]stakingprecompile.UnbondingDelegation)
				s.Require().EqualValues(len(unbonding), 1)
				// jsonData, _ := json.MarshalIndent(unbonding, "", "    ")
				// fmt.Printf("delegations: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)
			_, err = s.stakingKeeper.Undelegate(s.Ctx, delAddr, operatorAddress, sdk.NewDec(1))
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestDelegation() {
	method := stakingprecompile.StakingFunctionDelegation

	testCases := []struct {
		name        string
		malleate    func(delAddr, valAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr, valAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					valAddr,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				d := out[0].(stakingprecompile.Delegation)
				b := out[1].(*big.Int)
				_ = d
				_ = b
				/*
					jsonData, _ := json.MarshalIndent(d, "", "    ")
					fmt.Printf("delegation: %s\n", string(jsonData))
					fmt.Printf("balance: %v\n", b)
				*/
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes()), common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestUnbondingDelegation() {
	method := stakingprecompile.StakingFunctionUnbondingDelegation

	testCases := []struct {
		name        string
		malleate    func(delAddr, valAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr, valAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					valAddr,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				u := out[0].(stakingprecompile.UnbondingDelegation)
				_ = u
				// jsonData, _ := json.MarshalIndent(u, "", "    ")
				// fmt.Printf("delegation: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)
			_, err = s.stakingKeeper.Undelegate(s.Ctx, delAddr, operatorAddress, sdk.NewDec(1))
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes()), common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestDelegatorDelegations() {
	method := stakingprecompile.StakingFunctionDelegatorDelegations

	testCases := []struct {
		name        string
		malleate    func(delAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				d := out[0].([]stakingprecompile.DelegationResponse)
				paginationResult := out[1].(stakingprecompile.PageResponse)
				s.Assert().EqualValues(1, len(d))
				s.Assert().EqualValues(1, paginationResult.Total)
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("delegation: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestDelegatorUnbondingDelegations() {
	method := stakingprecompile.StakingFunctionDelegatorUnbondingDelegations

	testCases := []struct {
		name        string
		malleate    func(delAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				d := out[0].([]stakingprecompile.UnbondingDelegation)
				paginationResult := out[1].(stakingprecompile.PageResponse)
				s.Assert().EqualValues(1, len(d))
				s.Assert().EqualValues(1, paginationResult.Total)
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("delegation: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)
			_, err = s.stakingKeeper.Undelegate(s.Ctx, delAddr, operatorAddress, sdk.NewDec(1))
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestRedelegations() {
	method := stakingprecompile.StakingFunctionRedelegations

	testCases := []struct {
		name        string
		malleate    func(delAddr, srcValAddr, dstValAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr, srcValAddr, dstValAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					srcValAddr,
					dstValAddr,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				d := out[0].([]stakingprecompile.RedelegationResponse)
				paginationResult := out[1].(stakingprecompile.PageResponse)
				s.Assert().EqualValues(1, len(d))
				s.Assert().EqualValues(1, paginationResult.Total)
				// jsonData, _ := json.MarshalIndent(d, "", "    ")
				// fmt.Printf("redelegations: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)
			// setup redelegations
			s.setupValidator(s.signerOne)
			_, err = s.stakingKeeper.BeginRedelegation(s.Ctx, delAddr, operatorAddress, s.signerOne.ValAddr, sdk.NewDec(1))
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes()), common.Address(operatorAddress.Bytes()), common.Address(s.signerOne.ValAddr.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestDelegatorValidators() {
	method := stakingprecompile.StakingFunctionDelegatorValidators

	testCases := []struct {
		name        string
		malleate    func(delAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					query.PageRequest{
						Limit:      10,
						CountTotal: true,
					},
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				v := out[0].([]stakingprecompile.Validator)
				paginationResult := out[1].(stakingprecompile.PageResponse)
				s.Assert().EqualValues(1, len(v))
				s.Assert().EqualValues(1, paginationResult.Total)
				// jsonData, _ := json.MarshalIndent(v, "", "    ")
				// fmt.Printf("validators: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestDelegatorValidator() {
	method := stakingprecompile.StakingFunctionDelegatorValidator

	testCases := []struct {
		name        string
		malleate    func(delAddr, valAddr common.Address) []byte
		postCheck   func(bz []byte)
		gas         uint64
		expErr      bool
		errContains string
	}{
		{
			"success",
			func(delAddr, valAddr common.Address) []byte {
				input, err := s.abi.Pack(
					method,
					delAddr,
					valAddr,
				)
				s.Assert().NoError(err)
				return input
			},
			func(data []byte) {
				out, err := s.abi.Methods[method].Outputs.Unpack(data)
				s.Require().NoError(err, "failed to unpack output")
				v := out[0].(stakingprecompile.Validator)
				_ = v
				// jsonData, _ := json.MarshalIndent(v, "", "    ")
				// fmt.Printf("validators: %s\n", string(jsonData))
			},
			100000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			operatorAddress, err := s.firstBondedValidator()
			s.Require().NoError(err)
			d := s.stakingKeeper.GetValidatorDelegations(s.Ctx, operatorAddress)
			delAddr, err := sdk.AccAddressFromBech32(d[0].DelegatorAddress)
			s.Require().NoError(err)

			bz, err := s.runTx(tc.malleate(common.Address(delAddr.Bytes()), common.Address(operatorAddress.Bytes())), s.signerOne, 10000000)

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

func (s *StakingTestSuite) TestPool() {
	method := stakingprecompile.StakingFunctionPool

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
				bonded := out[0].(*big.Int)
				unbonded := out[0].(*big.Int)
				s.Assert().Equal(bonded.Int64(), int64(0))
				s.Assert().Equal(unbonded.Int64(), int64(0))
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

func (s *StakingTestSuite) TestParams() {
	method := stakingprecompile.StakingFunctionParams

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
				params := out[0].(stakingprecompile.Params)
				_ = params
				// jsonData, _ := json.MarshalIndent(params, "", "    ")
				// fmt.Printf("params: %s\n", string(jsonData))
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

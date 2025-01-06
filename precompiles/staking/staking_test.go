package staking_test

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"cosmossdk.io/math"
	stakingprecompile "github.com/0glabs/0g-chain/precompiles/staking"
	"github.com/0glabs/0g-chain/precompiles/testutil"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
)

type StakingTestSuite struct {
	testutil.PrecompileTestSuite

	abi           abi.ABI
	addr          common.Address
	staking       *stakingprecompile.StakingPrecompile
	stakingKeeper *stakingkeeper.Keeper
	signerOne     *testutil.TestSigner
	signerTwo     *testutil.TestSigner
}

func (suite *StakingTestSuite) SetupTest() {
	suite.PrecompileTestSuite.SetupTest()

	suite.stakingKeeper = suite.App.GetStakingKeeper()

	suite.addr = common.HexToAddress(stakingprecompile.PrecompileAddress)

	precompiles := suite.EvmKeeper.GetPrecompiles()
	precompile, ok := precompiles[suite.addr]
	suite.Assert().EqualValues(ok, true)

	suite.staking = precompile.(*stakingprecompile.StakingPrecompile)

	suite.signerOne = suite.GenSigner()
	suite.signerTwo = suite.GenSigner()

	abi, err := abi.JSON(strings.NewReader(stakingprecompile.StakingABI))
	suite.Assert().NoError(err)
	suite.abi = abi
}

func (suite *StakingTestSuite) AddDelegation(from string, to string, amount math.Int) {
	accAddr, err := sdk.AccAddressFromHexUnsafe(from)
	suite.Require().NoError(err)
	valAddr, err := sdk.ValAddressFromHex(to)
	suite.Require().NoError(err)
	validator, found := suite.StakingKeeper.GetValidator(suite.Ctx, valAddr)
	if !found {
		consPriv := ed25519.GenPrivKey()
		newValidator, err := stakingtypes.NewValidator(valAddr, consPriv.PubKey(), stakingtypes.Description{})
		suite.Require().NoError(err)
		validator = newValidator
	}
	validator.Tokens = validator.Tokens.Add(amount)
	validator.DelegatorShares = validator.DelegatorShares.Add(amount.ToLegacyDec())
	suite.StakingKeeper.SetValidator(suite.Ctx, validator)
	bonded := suite.stakingKeeper.GetDelegatorBonded(suite.Ctx, accAddr)
	suite.StakingKeeper.SetDelegation(suite.Ctx, stakingtypes.Delegation{
		DelegatorAddress: accAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           bonded.Add(amount).ToLegacyDec(),
	})
}

func (suite *StakingTestSuite) setupValidator(signer *testutil.TestSigner) {
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
	input, err := suite.abi.Pack(
		method,
		description,
		commission,
		minSelfDelegation,
		pubkey,
		value,
	)
	suite.Assert().NoError(err)
	_, err = suite.runTx(input, signer, 10000000)
	suite.Assert().NoError(err)
	_, err = suite.stakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.Ctx)
	suite.Assert().NoError(err)
}

func (suite *StakingTestSuite) firstBondedValidator() (sdk.ValAddress, error) {
	validators := suite.stakingKeeper.GetValidators(suite.Ctx, 10)
	for _, v := range validators {
		if v.IsBonded() {
			return sdk.ValAddressFromBech32(v.OperatorAddress)
		}
	}
	return nil, errors.New("no bonded validator")
}

func (suite *StakingTestSuite) runTx(input []byte, signer *testutil.TestSigner, gas uint64) ([]byte, error) {
	contract := vm.NewPrecompile(vm.AccountRef(signer.Addr), vm.AccountRef(suite.addr), big.NewInt(0), gas)
	contract.Input = input

	msgEthereumTx := evmtypes.NewTx(suite.EvmKeeper.ChainID(), 0, &suite.addr, big.NewInt(0), gas, big.NewInt(0), big.NewInt(0), big.NewInt(0), input, nil)
	msgEthereumTx.From = signer.HexAddr
	err := msgEthereumTx.Sign(suite.EthSigner, signer.Signer)
	suite.Assert().NoError(err, "failed to sign Ethereum message")

	proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
	cfg, err := suite.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.EvmKeeper.ChainID())
	suite.Assert().NoError(err, "failed to instantiate EVM config")

	msg, err := msgEthereumTx.AsMessage(suite.EthSigner, big.NewInt(0))
	suite.Assert().NoError(err, "failed to instantiate Ethereum message")

	evm := suite.EvmKeeper.NewEVM(suite.Ctx, msg, cfg, nil, suite.Statedb)
	precompiles := suite.EvmKeeper.GetPrecompiles()
	evm.WithPrecompiles(precompiles, []common.Address{suite.addr})

	bz, err := suite.staking.Run(evm, contract, false)
	if err == nil {
		evm.StateDB.(*statedb.StateDB).Commit()
	}
	return bz, err
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(StakingTestSuite))
}

package wrappeda0gibase_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/0glabs/0g-chain/precompiles/testutil"
	wrappeda0gibaseprecompile "github.com/0glabs/0g-chain/precompiles/wrapped-a0gi-base"
	wrappeda0gibasekeeper "github.com/0glabs/0g-chain/x/wrapped-a0gi-base/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
)

type WrappedA0giBaseTestSuite struct {
	testutil.PrecompileTestSuite

	abi             abi.ABI
	addr            common.Address
	wa0gibasekeeper wrappeda0gibasekeeper.Keeper
	wa0gibase       *wrappeda0gibaseprecompile.WrappedA0giBasePrecompile
	signerOne       *testutil.TestSigner
	signerTwo       *testutil.TestSigner
}

func (suite *WrappedA0giBaseTestSuite) SetupTest() {
	suite.PrecompileTestSuite.SetupTest()

	suite.wa0gibasekeeper = suite.App.GetWrappedA0GIBaseKeeper()

	suite.addr = common.HexToAddress(wrappeda0gibaseprecompile.PrecompileAddress)

	precompiles := suite.EvmKeeper.GetPrecompiles()
	precompile, ok := precompiles[suite.addr]
	suite.Assert().EqualValues(ok, true)

	suite.wa0gibase = precompile.(*wrappeda0gibaseprecompile.WrappedA0giBasePrecompile)

	suite.signerOne = suite.GenSigner()
	suite.signerTwo = suite.GenSigner()

	abi, err := abi.JSON(strings.NewReader(wrappeda0gibaseprecompile.Wrappeda0gibaseABI))
	suite.Assert().NoError(err)
	suite.abi = abi
}

func (suite *WrappedA0giBaseTestSuite) runTx(input []byte, signer *testutil.TestSigner, gas uint64) ([]byte, error) {
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

	bz, err := suite.wa0gibase.Run(evm, contract, false)
	if err == nil {
		evm.StateDB.(*statedb.StateDB).Commit()
	}
	return bz, err
}

func TestWrappedA0giBaseTestSuite(t *testing.T) {
	suite.Run(t, new(WrappedA0giBaseTestSuite))
}

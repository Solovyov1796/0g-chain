package staking

import (
	"errors"
	"math/big"

	precompiles_common "github.com/0glabs/0g-chain/precompiles/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

func (s *StakingPrecompile) CreateValidator(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgCreateValidator(args, contract.CallerAddress, s.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}
	// validation
	if contract.CallerAddress != evm.Origin {
		return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
	}
	// execute
	_, err = stakingkeeper.NewMsgServerImpl(s.stakingKeeper).CreateValidator(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack()
}

func (s *StakingPrecompile) EditValidator(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgEditValidator(args, contract.CallerAddress)
	if err != nil {
		return nil, err
	}
	// validation
	if contract.CallerAddress != evm.Origin {
		return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
	}
	// execute
	_, err = stakingkeeper.NewMsgServerImpl(s.stakingKeeper).EditValidator(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack()
}

func (s *StakingPrecompile) Delegate(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgDelegate(args, contract.CallerAddress, s.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}
	// validation
	/*
		if contract.CallerAddress != evm.Origin {
			return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
		}
	*/
	// execute
	_, err = stakingkeeper.NewMsgServerImpl(s.stakingKeeper).Delegate(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack()
}

func (s *StakingPrecompile) BeginRedelegate(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgBeginRedelegate(args, contract.CallerAddress, s.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}
	// validation
	/*
		if contract.CallerAddress != evm.Origin {
			return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
		}
	*/
	// execute
	response, err := stakingkeeper.NewMsgServerImpl(s.stakingKeeper).BeginRedelegate(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack(big.NewInt(response.CompletionTime.UTC().Unix()))
}

func (s *StakingPrecompile) Undelegate(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgUndelegate(args, contract.CallerAddress, s.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}
	// validation
	/*
		if contract.CallerAddress != evm.Origin {
			return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
		}
	*/
	// execute
	response, err := stakingkeeper.NewMsgServerImpl(s.stakingKeeper).Undelegate(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack(big.NewInt(response.CompletionTime.UTC().Unix()))
}

func (s *StakingPrecompile) CancelUnbondingDelegation(
	ctx sdk.Context,
	evm *vm.EVM,
	stateDB *statedb.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgCancelUnbondingDelegation(args, contract.CallerAddress, s.stakingKeeper.BondDenom(ctx))
	if err != nil {
		return nil, err
	}
	// validation
	/*
		if contract.CallerAddress != evm.Origin {
			return nil, errors.New(precompiles_common.ErrSenderNotOrigin)
		}
	*/
	// execute
	_, err = stakingkeeper.NewMsgServerImpl(s.stakingKeeper).CancelUnbondingDelegation(ctx, msg)
	if err != nil {
		return nil, err
	}
	// emit events
	return method.Outputs.Pack()
}

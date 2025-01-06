package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
)

func (s *StakingPrecompile) Validators(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryValidatorsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Validators(ctx, req)
	if err != nil {
		return nil, err
	}

	validators := make([]Validator, len(response.Validators))
	for i, v := range response.Validators {
		validators[i], err = convertValidator(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(validators, paginationResult)
}

func (s *StakingPrecompile) Validator(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryValidatorRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Validator(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := convertValidator(response.Validator)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res)
}

func (s *StakingPrecompile) ValidatorDelegations(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryValidatorDelegationsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.ValidatorDelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	delegationResponses := make([]DelegationResponse, len(response.DelegationResponses))
	for i, v := range response.DelegationResponses {
		delegationResponses[i], err = convertDelegationResponse(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(delegationResponses, paginationResult)
}

func (s *StakingPrecompile) ValidatorUnbondingDelegations(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryValidatorUnbondingDelegationsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.ValidatorUnbondingDelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	unbondingResponses := make([]UnbondingDelegation, len(response.UnbondingResponses))
	for i, v := range response.UnbondingResponses {
		unbondingResponses[i], err = convertUnbondingDelegation(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(unbondingResponses, paginationResult)
}

func (s *StakingPrecompile) Delegation(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryDelegationRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Delegation(ctx, req)
	if err != nil {
		return nil, err
	}
	delegation, err := convertDelegation(response.DelegationResponse.Delegation)
	if err != nil {
		return nil, err
	}
	balance := response.DelegationResponse.Balance.Amount.BigInt()

	return method.Outputs.Pack(delegation, balance)
}

func (s *StakingPrecompile) UnbondingDelegation(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryUnbondingDelegationRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.UnbondingDelegation(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := convertUnbondingDelegation(response.Unbond)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res)
}

func (s *StakingPrecompile) DelegatorDelegations(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryDelegatorDelegationsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.DelegatorDelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	delegationResponses := make([]DelegationResponse, len(response.DelegationResponses))
	for i, v := range response.DelegationResponses {
		delegationResponses[i], err = convertDelegationResponse(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(delegationResponses, paginationResult)
}

func (s *StakingPrecompile) DelegatorUnbondingDelegations(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryDelegatorUnbondingDelegationsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.DelegatorUnbondingDelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	unbondingResponses := make([]UnbondingDelegation, len(response.UnbondingResponses))
	for i, v := range response.UnbondingResponses {
		unbondingResponses[i], err = convertUnbondingDelegation(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(unbondingResponses, paginationResult)
}

func (s *StakingPrecompile) Redelegations(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryRedelegationsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Redelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	redelegationResponses := make([]RedelegationResponse, len(response.RedelegationResponses))
	for i, v := range response.RedelegationResponses {
		redelegationResponses[i], err = convertRedelegationResponse(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(redelegationResponses, paginationResult)
}

func (s *StakingPrecompile) DelegatorValidators(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryDelegatorValidatorsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.DelegatorValidators(ctx, req)
	if err != nil {
		return nil, err
	}

	validators := make([]Validator, len(response.Validators))
	for i, v := range response.Validators {
		validators[i], err = convertValidator(v)
		if err != nil {
			return nil, err
		}
	}
	paginationResult := convertPageResponse(response.Pagination)

	return method.Outputs.Pack(validators, paginationResult)
}

func (s *StakingPrecompile) DelegatorValidator(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryDelegatorValidatorRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.DelegatorValidator(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := convertValidator(response.Validator)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res)
}

func (s *StakingPrecompile) Pool(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryPoolRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Pool(ctx, req)
	if err != nil {
		return nil, err
	}
	notBondedTokens := response.Pool.NotBondedTokens.BigInt()
	bondedTokens := response.Pool.BondedTokens.BigInt()

	return method.Outputs.Pack(notBondedTokens, bondedTokens)
}

func (s *StakingPrecompile) Params(ctx sdk.Context, _ *vm.EVM, method *abi.Method, args []interface{}) ([]byte, error) {
	req, err := NewQueryParamsRequest(args)
	if err != nil {
		return nil, err
	}
	response, err := stakingkeeper.Querier{Keeper: s.stakingKeeper}.Params(ctx, req)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(convertParams(response.Params))
}

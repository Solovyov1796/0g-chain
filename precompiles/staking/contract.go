// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorSrcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorDstAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"beginRedelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"creationHeight\",\"type\":\"uint256\"}],\"name\":\"cancelUnbondingDelegation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commission\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"pubkey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"createValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"internalType\":\"structDelegation\",\"name\":\"delegation\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"delegatorDelegations\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"internalType\":\"structDelegation\",\"name\":\"delegation\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structDelegationResponse[]\",\"name\":\"delegationResponses\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"delegatorUnbondingDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structUnbondingDelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structUnbondingDelegation[]\",\"name\":\"unbondingResponses\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"}],\"name\":\"delegatorValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commissionRates\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"updateTime\",\"type\":\"uint256\"}],\"internalType\":\"structCommission\",\"name\":\"commission\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"},{\"internalType\":\"uint64[]\",\"name\":\"unbondingIds\",\"type\":\"uint64[]\"}],\"internalType\":\"structValidator\",\"name\":\"validator\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"delegatorValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commissionRates\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"updateTime\",\"type\":\"uint256\"}],\"internalType\":\"structCommission\",\"name\":\"commission\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"},{\"internalType\":\"uint64[]\",\"name\":\"unbondingIds\",\"type\":\"uint64[]\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isNull\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structNullableUint\",\"name\":\"commissionRate\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isNull\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structNullableUint\",\"name\":\"minSelfDelegation\",\"type\":\"tuple\"}],\"name\":\"editValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"params\",\"outputs\":[{\"components\":[{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"maxValidators\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"maxEntries\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"historicalEntries\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"bondDenom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"minCommissionRate\",\"type\":\"uint256\"}],\"internalType\":\"structParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"notBondedTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedTokens\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"srcValidatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dstValidatorAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pageRequest\",\"type\":\"tuple\"}],\"name\":\"redelegations\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorSrcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorDstAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sharesDst\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structRedelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structRedelegation\",\"name\":\"redelegation\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sharesDst\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structRedelegationEntry\",\"name\":\"redelegationEntry\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structRedelegationEntryResponse[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structRedelegationResponse[]\",\"name\":\"redelegationResponses\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"}],\"name\":\"unbondingDelegation\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structUnbondingDelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structUnbondingDelegation\",\"name\":\"unbond\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"validator\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commissionRates\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"updateTime\",\"type\":\"uint256\"}],\"internalType\":\"structCommission\",\"name\":\"commission\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"},{\"internalType\":\"uint64[]\",\"name\":\"unbondingIds\",\"type\":\"uint64[]\"}],\"internalType\":\"structValidator\",\"name\":\"validator\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"validatorDelegations\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"internalType\":\"structDelegation\",\"name\":\"delegation\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structDelegationResponse[]\",\"name\":\"delegationResponses\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddr\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"validatorUnbondingDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structUnbondingDelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structUnbondingDelegation[]\",\"name\":\"unbondingResponses\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"status\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pagination\",\"type\":\"tuple\"}],\"name\":\"validators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commissionRates\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"updateTime\",\"type\":\"uint256\"}],\"internalType\":\"structCommission\",\"name\":\"commission\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"},{\"internalType\":\"uint64[]\",\"name\":\"unbondingIds\",\"type\":\"uint64[]\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"paginationResult\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// Delegation is a free data retrieval call binding the contract method 0x046d3307.
//
// Solidity: function delegation(address delegatorAddr, address validatorAddr) view returns((address,address,uint256) delegation, uint256 balance)
func (_Staking *StakingCaller) Delegation(opts *bind.CallOpts, delegatorAddr common.Address, validatorAddr common.Address) (struct {
	Delegation Delegation
	Balance    *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegation", delegatorAddr, validatorAddr)

	outstruct := new(struct {
		Delegation Delegation
		Balance    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Delegation = *abi.ConvertType(out[0], new(Delegation)).(*Delegation)
	outstruct.Balance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Delegation is a free data retrieval call binding the contract method 0x046d3307.
//
// Solidity: function delegation(address delegatorAddr, address validatorAddr) view returns((address,address,uint256) delegation, uint256 balance)
func (_Staking *StakingSession) Delegation(delegatorAddr common.Address, validatorAddr common.Address) (struct {
	Delegation Delegation
	Balance    *big.Int
}, error) {
	return _Staking.Contract.Delegation(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// Delegation is a free data retrieval call binding the contract method 0x046d3307.
//
// Solidity: function delegation(address delegatorAddr, address validatorAddr) view returns((address,address,uint256) delegation, uint256 balance)
func (_Staking *StakingCallerSession) Delegation(delegatorAddr common.Address, validatorAddr common.Address) (struct {
	Delegation Delegation
	Balance    *big.Int
}, error) {
	return _Staking.Contract.Delegation(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// DelegatorDelegations is a free data retrieval call binding the contract method 0x256bd907.
//
// Solidity: function delegatorDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) DelegatorDelegations(opts *bind.CallOpts, delegatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegatorDelegations", delegatorAddr, pagination)

	outstruct := new(struct {
		DelegationResponses []DelegationResponse
		PaginationResult    PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DelegationResponses = *abi.ConvertType(out[0], new([]DelegationResponse)).(*[]DelegationResponse)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// DelegatorDelegations is a free data retrieval call binding the contract method 0x256bd907.
//
// Solidity: function delegatorDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) DelegatorDelegations(delegatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	return _Staking.Contract.DelegatorDelegations(&_Staking.CallOpts, delegatorAddr, pagination)
}

// DelegatorDelegations is a free data retrieval call binding the contract method 0x256bd907.
//
// Solidity: function delegatorDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) DelegatorDelegations(delegatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	return _Staking.Contract.DelegatorDelegations(&_Staking.CallOpts, delegatorAddr, pagination)
}

// DelegatorUnbondingDelegations is a free data retrieval call binding the contract method 0x155fd3ff.
//
// Solidity: function delegatorUnbondingDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) DelegatorUnbondingDelegations(opts *bind.CallOpts, delegatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegatorUnbondingDelegations", delegatorAddr, pagination)

	outstruct := new(struct {
		UnbondingResponses []UnbondingDelegation
		PaginationResult   PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UnbondingResponses = *abi.ConvertType(out[0], new([]UnbondingDelegation)).(*[]UnbondingDelegation)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// DelegatorUnbondingDelegations is a free data retrieval call binding the contract method 0x155fd3ff.
//
// Solidity: function delegatorUnbondingDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) DelegatorUnbondingDelegations(delegatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	return _Staking.Contract.DelegatorUnbondingDelegations(&_Staking.CallOpts, delegatorAddr, pagination)
}

// DelegatorUnbondingDelegations is a free data retrieval call binding the contract method 0x155fd3ff.
//
// Solidity: function delegatorUnbondingDelegations(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) DelegatorUnbondingDelegations(delegatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	return _Staking.Contract.DelegatorUnbondingDelegations(&_Staking.CallOpts, delegatorAddr, pagination)
}

// DelegatorValidator is a free data retrieval call binding the contract method 0x77007d17.
//
// Solidity: function delegatorValidator(address delegatorAddr, address validatorAddr) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingCaller) DelegatorValidator(opts *bind.CallOpts, delegatorAddr common.Address, validatorAddr common.Address) (Validator, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegatorValidator", delegatorAddr, validatorAddr)

	if err != nil {
		return *new(Validator), err
	}

	out0 := *abi.ConvertType(out[0], new(Validator)).(*Validator)

	return out0, err

}

// DelegatorValidator is a free data retrieval call binding the contract method 0x77007d17.
//
// Solidity: function delegatorValidator(address delegatorAddr, address validatorAddr) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingSession) DelegatorValidator(delegatorAddr common.Address, validatorAddr common.Address) (Validator, error) {
	return _Staking.Contract.DelegatorValidator(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// DelegatorValidator is a free data retrieval call binding the contract method 0x77007d17.
//
// Solidity: function delegatorValidator(address delegatorAddr, address validatorAddr) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingCallerSession) DelegatorValidator(delegatorAddr common.Address, validatorAddr common.Address) (Validator, error) {
	return _Staking.Contract.DelegatorValidator(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// DelegatorValidators is a free data retrieval call binding the contract method 0xe8e4fe99.
//
// Solidity: function delegatorValidators(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) DelegatorValidators(opts *bind.CallOpts, delegatorAddr common.Address, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegatorValidators", delegatorAddr, pagination)

	outstruct := new(struct {
		Validators       []Validator
		PaginationResult PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validators = *abi.ConvertType(out[0], new([]Validator)).(*[]Validator)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// DelegatorValidators is a free data retrieval call binding the contract method 0xe8e4fe99.
//
// Solidity: function delegatorValidators(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) DelegatorValidators(delegatorAddr common.Address, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	return _Staking.Contract.DelegatorValidators(&_Staking.CallOpts, delegatorAddr, pagination)
}

// DelegatorValidators is a free data retrieval call binding the contract method 0xe8e4fe99.
//
// Solidity: function delegatorValidators(address delegatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) DelegatorValidators(delegatorAddr common.Address, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	return _Staking.Contract.DelegatorValidators(&_Staking.CallOpts, delegatorAddr, pagination)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((int64,uint32,uint32,uint32,string,uint256) params)
func (_Staking *StakingCaller) Params(opts *bind.CallOpts) (Params, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "params")

	if err != nil {
		return *new(Params), err
	}

	out0 := *abi.ConvertType(out[0], new(Params)).(*Params)

	return out0, err

}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((int64,uint32,uint32,uint32,string,uint256) params)
func (_Staking *StakingSession) Params() (Params, error) {
	return _Staking.Contract.Params(&_Staking.CallOpts)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((int64,uint32,uint32,uint32,string,uint256) params)
func (_Staking *StakingCallerSession) Params() (Params, error) {
	return _Staking.Contract.Params(&_Staking.CallOpts)
}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(uint256 notBondedTokens, uint256 bondedTokens)
func (_Staking *StakingCaller) Pool(opts *bind.CallOpts) (struct {
	NotBondedTokens *big.Int
	BondedTokens    *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "pool")

	outstruct := new(struct {
		NotBondedTokens *big.Int
		BondedTokens    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NotBondedTokens = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BondedTokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(uint256 notBondedTokens, uint256 bondedTokens)
func (_Staking *StakingSession) Pool() (struct {
	NotBondedTokens *big.Int
	BondedTokens    *big.Int
}, error) {
	return _Staking.Contract.Pool(&_Staking.CallOpts)
}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(uint256 notBondedTokens, uint256 bondedTokens)
func (_Staking *StakingCallerSession) Pool() (struct {
	NotBondedTokens *big.Int
	BondedTokens    *big.Int
}, error) {
	return _Staking.Contract.Pool(&_Staking.CallOpts)
}

// Redelegations is a free data retrieval call binding the contract method 0xeb5643a9.
//
// Solidity: function redelegations(address delegatorAddress, address srcValidatorAddress, address dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((address,address,address,(int64,int64,uint256,uint256,uint64,int64)[]),((int64,int64,uint256,uint256,uint64,int64),uint256)[])[] redelegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) Redelegations(opts *bind.CallOpts, delegatorAddress common.Address, srcValidatorAddress common.Address, dstValidatorAddress common.Address, pageRequest PageRequest) (struct {
	RedelegationResponses []RedelegationResponse
	PaginationResult      PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "redelegations", delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)

	outstruct := new(struct {
		RedelegationResponses []RedelegationResponse
		PaginationResult      PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RedelegationResponses = *abi.ConvertType(out[0], new([]RedelegationResponse)).(*[]RedelegationResponse)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// Redelegations is a free data retrieval call binding the contract method 0xeb5643a9.
//
// Solidity: function redelegations(address delegatorAddress, address srcValidatorAddress, address dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((address,address,address,(int64,int64,uint256,uint256,uint64,int64)[]),((int64,int64,uint256,uint256,uint64,int64),uint256)[])[] redelegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) Redelegations(delegatorAddress common.Address, srcValidatorAddress common.Address, dstValidatorAddress common.Address, pageRequest PageRequest) (struct {
	RedelegationResponses []RedelegationResponse
	PaginationResult      PageResponse
}, error) {
	return _Staking.Contract.Redelegations(&_Staking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)
}

// Redelegations is a free data retrieval call binding the contract method 0xeb5643a9.
//
// Solidity: function redelegations(address delegatorAddress, address srcValidatorAddress, address dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((address,address,address,(int64,int64,uint256,uint256,uint64,int64)[]),((int64,int64,uint256,uint256,uint64,int64),uint256)[])[] redelegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) Redelegations(delegatorAddress common.Address, srcValidatorAddress common.Address, dstValidatorAddress common.Address, pageRequest PageRequest) (struct {
	RedelegationResponses []RedelegationResponse
	PaginationResult      PageResponse
}, error) {
	return _Staking.Contract.Redelegations(&_Staking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)
}

// UnbondingDelegation is a free data retrieval call binding the contract method 0x97e41907.
//
// Solidity: function unbondingDelegation(address delegatorAddr, address validatorAddr) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[]) unbond)
func (_Staking *StakingCaller) UnbondingDelegation(opts *bind.CallOpts, delegatorAddr common.Address, validatorAddr common.Address) (UnbondingDelegation, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "unbondingDelegation", delegatorAddr, validatorAddr)

	if err != nil {
		return *new(UnbondingDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new(UnbondingDelegation)).(*UnbondingDelegation)

	return out0, err

}

// UnbondingDelegation is a free data retrieval call binding the contract method 0x97e41907.
//
// Solidity: function unbondingDelegation(address delegatorAddr, address validatorAddr) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[]) unbond)
func (_Staking *StakingSession) UnbondingDelegation(delegatorAddr common.Address, validatorAddr common.Address) (UnbondingDelegation, error) {
	return _Staking.Contract.UnbondingDelegation(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// UnbondingDelegation is a free data retrieval call binding the contract method 0x97e41907.
//
// Solidity: function unbondingDelegation(address delegatorAddr, address validatorAddr) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[]) unbond)
func (_Staking *StakingCallerSession) UnbondingDelegation(delegatorAddr common.Address, validatorAddr common.Address) (UnbondingDelegation, error) {
	return _Staking.Contract.UnbondingDelegation(&_Staking.CallOpts, delegatorAddr, validatorAddr)
}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingCaller) Validator(opts *bind.CallOpts, validatorAddress common.Address) (Validator, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "validator", validatorAddress)

	if err != nil {
		return *new(Validator), err
	}

	out0 := *abi.ConvertType(out[0], new(Validator)).(*Validator)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingSession) Validator(validatorAddress common.Address) (Validator, error) {
	return _Staking.Contract.Validator(&_Staking.CallOpts, validatorAddress)
}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[]) validator)
func (_Staking *StakingCallerSession) Validator(validatorAddress common.Address) (Validator, error) {
	return _Staking.Contract.Validator(&_Staking.CallOpts, validatorAddress)
}

// ValidatorDelegations is a free data retrieval call binding the contract method 0xc89017cd.
//
// Solidity: function validatorDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) ValidatorDelegations(opts *bind.CallOpts, validatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "validatorDelegations", validatorAddr, pagination)

	outstruct := new(struct {
		DelegationResponses []DelegationResponse
		PaginationResult    PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DelegationResponses = *abi.ConvertType(out[0], new([]DelegationResponse)).(*[]DelegationResponse)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// ValidatorDelegations is a free data retrieval call binding the contract method 0xc89017cd.
//
// Solidity: function validatorDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) ValidatorDelegations(validatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	return _Staking.Contract.ValidatorDelegations(&_Staking.CallOpts, validatorAddr, pagination)
}

// ValidatorDelegations is a free data retrieval call binding the contract method 0xc89017cd.
//
// Solidity: function validatorDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns(((address,address,uint256),uint256)[] delegationResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) ValidatorDelegations(validatorAddr common.Address, pagination PageRequest) (struct {
	DelegationResponses []DelegationResponse
	PaginationResult    PageResponse
}, error) {
	return _Staking.Contract.ValidatorDelegations(&_Staking.CallOpts, validatorAddr, pagination)
}

// ValidatorUnbondingDelegations is a free data retrieval call binding the contract method 0x51b55195.
//
// Solidity: function validatorUnbondingDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) ValidatorUnbondingDelegations(opts *bind.CallOpts, validatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "validatorUnbondingDelegations", validatorAddr, pagination)

	outstruct := new(struct {
		UnbondingResponses []UnbondingDelegation
		PaginationResult   PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UnbondingResponses = *abi.ConvertType(out[0], new([]UnbondingDelegation)).(*[]UnbondingDelegation)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// ValidatorUnbondingDelegations is a free data retrieval call binding the contract method 0x51b55195.
//
// Solidity: function validatorUnbondingDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) ValidatorUnbondingDelegations(validatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	return _Staking.Contract.ValidatorUnbondingDelegations(&_Staking.CallOpts, validatorAddr, pagination)
}

// ValidatorUnbondingDelegations is a free data retrieval call binding the contract method 0x51b55195.
//
// Solidity: function validatorUnbondingDelegations(address validatorAddr, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,address,(int64,int64,uint256,uint256,uint64,int64)[])[] unbondingResponses, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) ValidatorUnbondingDelegations(validatorAddr common.Address, pagination PageRequest) (struct {
	UnbondingResponses []UnbondingDelegation
	PaginationResult   PageResponse
}, error) {
	return _Staking.Contract.ValidatorUnbondingDelegations(&_Staking.CallOpts, validatorAddr, pagination)
}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingCaller) Validators(opts *bind.CallOpts, status string, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "validators", status, pagination)

	outstruct := new(struct {
		Validators       []Validator
		PaginationResult PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validators = *abi.ConvertType(out[0], new([]Validator)).(*[]Validator)
	outstruct.PaginationResult = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingSession) Validators(status string, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	return _Staking.Contract.Validators(&_Staking.CallOpts, status, pagination)
}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pagination) view returns((address,string,bool,uint8,uint256,uint256,(string,string,string,string,string),int64,int64,((uint256,uint256,uint256),uint256),uint256,int64,uint64[])[] validators, (bytes,uint64) paginationResult)
func (_Staking *StakingCallerSession) Validators(status string, pagination PageRequest) (struct {
	Validators       []Validator
	PaginationResult PageResponse
}, error) {
	return _Staking.Contract.Validators(&_Staking.CallOpts, status, pagination)
}

// BeginRedelegate is a paid mutator transaction binding the contract method 0xb3a8ae3b.
//
// Solidity: function beginRedelegate(address validatorSrcAddress, address validatorDstAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingTransactor) BeginRedelegate(opts *bind.TransactOpts, validatorSrcAddress common.Address, validatorDstAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "beginRedelegate", validatorSrcAddress, validatorDstAddress, amount)
}

// BeginRedelegate is a paid mutator transaction binding the contract method 0xb3a8ae3b.
//
// Solidity: function beginRedelegate(address validatorSrcAddress, address validatorDstAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingSession) BeginRedelegate(validatorSrcAddress common.Address, validatorDstAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.BeginRedelegate(&_Staking.TransactOpts, validatorSrcAddress, validatorDstAddress, amount)
}

// BeginRedelegate is a paid mutator transaction binding the contract method 0xb3a8ae3b.
//
// Solidity: function beginRedelegate(address validatorSrcAddress, address validatorDstAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingTransactorSession) BeginRedelegate(validatorSrcAddress common.Address, validatorDstAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.BeginRedelegate(&_Staking.TransactOpts, validatorSrcAddress, validatorDstAddress, amount)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x50826aef.
//
// Solidity: function cancelUnbondingDelegation(address validatorAddress, uint256 amount, uint256 creationHeight) returns()
func (_Staking *StakingTransactor) CancelUnbondingDelegation(opts *bind.TransactOpts, validatorAddress common.Address, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "cancelUnbondingDelegation", validatorAddress, amount, creationHeight)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x50826aef.
//
// Solidity: function cancelUnbondingDelegation(address validatorAddress, uint256 amount, uint256 creationHeight) returns()
func (_Staking *StakingSession) CancelUnbondingDelegation(validatorAddress common.Address, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CancelUnbondingDelegation(&_Staking.TransactOpts, validatorAddress, amount, creationHeight)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x50826aef.
//
// Solidity: function cancelUnbondingDelegation(address validatorAddress, uint256 amount, uint256 creationHeight) returns()
func (_Staking *StakingTransactorSession) CancelUnbondingDelegation(validatorAddress common.Address, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CancelUnbondingDelegation(&_Staking.TransactOpts, validatorAddress, amount, creationHeight)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x3adeba1e.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commission, uint256 minSelfDelegation, string pubkey, uint256 value) returns()
func (_Staking *StakingTransactor) CreateValidator(opts *bind.TransactOpts, description Description, commission CommissionRates, minSelfDelegation *big.Int, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "createValidator", description, commission, minSelfDelegation, pubkey, value)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x3adeba1e.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commission, uint256 minSelfDelegation, string pubkey, uint256 value) returns()
func (_Staking *StakingSession) CreateValidator(description Description, commission CommissionRates, minSelfDelegation *big.Int, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CreateValidator(&_Staking.TransactOpts, description, commission, minSelfDelegation, pubkey, value)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x3adeba1e.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commission, uint256 minSelfDelegation, string pubkey, uint256 value) returns()
func (_Staking *StakingTransactorSession) CreateValidator(description Description, commission CommissionRates, minSelfDelegation *big.Int, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CreateValidator(&_Staking.TransactOpts, description, commission, minSelfDelegation, pubkey, value)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactor) Delegate(opts *bind.TransactOpts, validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegate", validatorAddress, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorAddress, uint256 amount) returns()
func (_Staking *StakingSession) Delegate(validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validatorAddress, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactorSession) Delegate(validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validatorAddress, amount)
}

// EditValidator is a paid mutator transaction binding the contract method 0x34397fbb.
//
// Solidity: function editValidator((string,string,string,string,string) description, (bool,uint256) commissionRate, (bool,uint256) minSelfDelegation) returns()
func (_Staking *StakingTransactor) EditValidator(opts *bind.TransactOpts, description Description, commissionRate NullableUint, minSelfDelegation NullableUint) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "editValidator", description, commissionRate, minSelfDelegation)
}

// EditValidator is a paid mutator transaction binding the contract method 0x34397fbb.
//
// Solidity: function editValidator((string,string,string,string,string) description, (bool,uint256) commissionRate, (bool,uint256) minSelfDelegation) returns()
func (_Staking *StakingSession) EditValidator(description Description, commissionRate NullableUint, minSelfDelegation NullableUint) (*types.Transaction, error) {
	return _Staking.Contract.EditValidator(&_Staking.TransactOpts, description, commissionRate, minSelfDelegation)
}

// EditValidator is a paid mutator transaction binding the contract method 0x34397fbb.
//
// Solidity: function editValidator((string,string,string,string,string) description, (bool,uint256) commissionRate, (bool,uint256) minSelfDelegation) returns()
func (_Staking *StakingTransactorSession) EditValidator(description Description, commissionRate NullableUint, minSelfDelegation NullableUint) (*types.Transaction, error) {
	return _Staking.Contract.EditValidator(&_Staking.TransactOpts, description, commissionRate, minSelfDelegation)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validatorAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingTransactor) Undelegate(opts *bind.TransactOpts, validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegate", validatorAddress, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validatorAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingSession) Undelegate(validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validatorAddress, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validatorAddress, uint256 amount) returns(uint256 completionTime)
func (_Staking *StakingTransactorSession) Undelegate(validatorAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validatorAddress, amount)
}

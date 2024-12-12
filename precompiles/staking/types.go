package staking

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	precompiles_common "github.com/0glabs/0g-chain/precompiles/common"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

type Commission = struct {
	CommissionRates CommissionRates `json:"commissionRates"`
	UpdateTime      *big.Int        `json:"updateTime"`
}

type CommissionRates = struct {
	Rate          *big.Int `json:"rate"`
	MaxRate       *big.Int `json:"maxRate"`
	MaxChangeRate *big.Int `json:"maxChangeRate"`
}

type Delegation = struct {
	DelegatorAddress common.Address `json:"delegatorAddress"`
	ValidatorAddress common.Address `json:"validatorAddress"`
	Shares           *big.Int       `json:"shares"`
}

type DelegationResponse = struct {
	Delegation Delegation `json:"delegation"`
	Balance    *big.Int   `json:"balance"`
}

type Description = struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"securityContact"`
	Details         string `json:"details"`
}

type NullableUint = struct {
	IsNull bool     `json:"isNull"`
	Value  *big.Int `json:"value"`
}

type PageRequest = struct {
	Key        []byte `json:"key"`
	Offset     uint64 `json:"offset"`
	Limit      uint64 `json:"limit"`
	CountTotal bool   `json:"countTotal"`
	Reverse    bool   `json:"reverse"`
}

type PageResponse = struct {
	NextKey []byte `json:"nextKey"`
	Total   uint64 `json:"total"`
}

type Params = struct {
	UnbondingTime     int64    `json:"unbondingTime"`
	MaxValidators     uint32   `json:"maxValidators"`
	MaxEntries        uint32   `json:"maxEntries"`
	HistoricalEntries uint32   `json:"historicalEntries"`
	BondDenom         string   `json:"bondDenom"`
	MinCommissionRate *big.Int `json:"minCommissionRate"`
}

type Redelegation = struct {
	DelegatorAddress    common.Address      `json:"delegatorAddress"`
	ValidatorSrcAddress common.Address      `json:"validatorSrcAddress"`
	ValidatorDstAddress common.Address      `json:"validatorDstAddress"`
	Entries             []RedelegationEntry `json:"entries"`
}

type RedelegationEntry = struct {
	CreationHeight          int64    `json:"creationHeight"`
	CompletionTime          int64    `json:"completionTime"`
	InitialBalance          *big.Int `json:"initialBalance"`
	SharesDst               *big.Int `json:"sharesDst"`
	UnbondingId             uint64   `json:"unbondingId"`
	UnbondingOnHoldRefCount int64    `json:"unbondingOnHoldRefCount"`
}

type RedelegationEntryResponse = struct {
	RedelegationEntry RedelegationEntry `json:"redelegationEntry"`
	Balance           *big.Int          `json:"balance"`
}

type RedelegationResponse = struct {
	Redelegation Redelegation                `json:"redelegation"`
	Entries      []RedelegationEntryResponse `json:"entries"`
}

type UnbondingDelegation = struct {
	DelegatorAddress common.Address             `json:"delegatorAddress"`
	ValidatorAddress common.Address             `json:"validatorAddress"`
	Entries          []UnbondingDelegationEntry `json:"entries"`
}

type UnbondingDelegationEntry = struct {
	CreationHeight          int64    `json:"creationHeight"`
	CompletionTime          int64    `json:"completionTime"`
	InitialBalance          *big.Int `json:"initialBalance"`
	Balance                 *big.Int `json:"balance"`
	UnbondingId             uint64   `json:"unbondingId"`
	UnbondingOnHoldRefCount int64    `json:"unbondingOnHoldRefCount"`
}

type Validator = struct {
	OperatorAddress         common.Address `json:"operatorAddress"`
	ConsensusPubkey         string         `json:"consensusPubkey"`
	Jailed                  bool           `json:"jailed"`
	Status                  uint8          `json:"status"`
	Tokens                  *big.Int       `json:"tokens"`
	DelegatorShares         *big.Int       `json:"delegatorShares"`
	Description             Description    `json:"description"`
	UnbondingHeight         int64          `json:"unbondingHeight"`
	UnbondingTime           int64          `json:"unbondingTime"`
	Commission              Commission     `json:"commission"`
	MinSelfDelegation       *big.Int       `json:"minSelfDelegation"`
	UnbondingOnHoldRefCount int64          `json:"unbondingOnHoldRefCount"`
	UnbondingIds            []uint64       `json:"unbondingIds"`
}

func convertValidator(v stakingtypes.Validator) (Validator, error) {
	validator := Validator{}
	operatorAddress, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		return validator, err
	}
	validator.OperatorAddress = common.BytesToAddress(operatorAddress.Bytes())

	ed25519pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		validator.ConsensusPubkey = v.ConsensusPubkey.String()
	} else {
		validator.ConsensusPubkey = base64.StdEncoding.EncodeToString(ed25519pk.Bytes())
	}

	validator.Jailed = v.Jailed
	validator.Status = uint8(v.Status)
	validator.Tokens = v.Tokens.BigInt()
	validator.DelegatorShares = v.DelegatorShares.BigInt()
	validator.Description = Description{
		Moniker:         v.Description.Moniker,
		Identity:        v.Description.Identity,
		Website:         v.Description.Website,
		SecurityContact: v.Description.SecurityContact,
		Details:         v.Description.Details,
	}
	validator.UnbondingHeight = v.UnbondingHeight
	validator.UnbondingTime = v.UnbondingTime.UTC().Unix()
	validator.Commission = Commission{
		CommissionRates: convertCommissionRates(v.Commission.CommissionRates),
		UpdateTime:      big.NewInt(v.Commission.UpdateTime.UTC().Unix()),
	}
	validator.MinSelfDelegation = v.MinSelfDelegation.BigInt()
	validator.UnbondingOnHoldRefCount = v.UnbondingOnHoldRefCount
	validator.UnbondingIds = v.UnbondingIds
	return validator, nil
}

func convertQueryPageRequest(pagination PageRequest) *query.PageRequest {
	return &query.PageRequest{
		Key:        pagination.Key,
		Offset:     pagination.Offset,
		Limit:      pagination.Limit,
		CountTotal: pagination.CountTotal,
		Reverse:    pagination.Reverse,
	}
}

func convertPageResponse(pagination *query.PageResponse) PageResponse {
	if pagination == nil {
		return PageResponse{
			NextKey: make([]byte, 0),
			Total:   1,
		}
	}
	return PageResponse{
		NextKey: pagination.NextKey,
		Total:   pagination.Total,
	}
}

func convertStakingDescription(description Description) stakingtypes.Description {
	return stakingtypes.Description{
		Moniker:         description.Moniker,
		Identity:        description.Identity,
		Website:         description.Website,
		SecurityContact: description.SecurityContact,
		Details:         description.Details,
	}
}

func convertStakingCommissionRates(commission CommissionRates) stakingtypes.CommissionRates {
	return stakingtypes.CommissionRates{
		Rate:          precompiles_common.BigIntToLegacyDec(commission.Rate),
		MaxRate:       precompiles_common.BigIntToLegacyDec(commission.MaxRate),
		MaxChangeRate: precompiles_common.BigIntToLegacyDec(commission.MaxChangeRate),
	}
}

func convertCommissionRates(commission stakingtypes.CommissionRates) CommissionRates {
	return CommissionRates{
		Rate:          commission.Rate.BigInt(),
		MaxRate:       commission.MaxRate.BigInt(),
		MaxChangeRate: commission.MaxChangeRate.BigInt(),
	}
}

func convertDelegation(delegation stakingtypes.Delegation) (Delegation, error) {
	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return Delegation{}, err
	}
	validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		return Delegation{}, err
	}
	return Delegation{
		DelegatorAddress: common.BytesToAddress(delegatorAddress.Bytes()),
		ValidatorAddress: common.BytesToAddress(validatorAddress.Bytes()),
		Shares:           delegation.Shares.BigInt(),
	}, nil
}

func convertDelegationResponse(response stakingtypes.DelegationResponse) (DelegationResponse, error) {
	delegation, err := convertDelegation(response.Delegation)
	if err != nil {
		return DelegationResponse{}, err
	}
	return DelegationResponse{
		Delegation: delegation,
		Balance:    response.Balance.Amount.BigInt(),
	}, nil
}

func convertUnbondingDelegationEntry(entry stakingtypes.UnbondingDelegationEntry) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight:          entry.CreationHeight,
		CompletionTime:          entry.CompletionTime.UTC().Unix(),
		InitialBalance:          entry.InitialBalance.BigInt(),
		Balance:                 entry.Balance.BigInt(),
		UnbondingId:             entry.UnbondingId,
		UnbondingOnHoldRefCount: entry.UnbondingOnHoldRefCount,
	}
}

func convertUnbondingDelegation(response stakingtypes.UnbondingDelegation) (UnbondingDelegation, error) {
	entries := make([]UnbondingDelegationEntry, len(response.Entries))
	for i, v := range response.Entries {
		entries[i] = convertUnbondingDelegationEntry(v)
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(response.DelegatorAddress)
	if err != nil {
		return UnbondingDelegation{}, err
	}
	validatorAddress, err := sdk.ValAddressFromBech32(response.ValidatorAddress)
	if err != nil {
		return UnbondingDelegation{}, err
	}
	return UnbondingDelegation{
		DelegatorAddress: common.BytesToAddress(delegatorAddress.Bytes()),
		ValidatorAddress: common.BytesToAddress(validatorAddress.Bytes()),
		Entries:          entries,
	}, nil
}

func convertRedelegationEntry(entry stakingtypes.RedelegationEntry) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight:          entry.CreationHeight,
		CompletionTime:          entry.CompletionTime.UTC().Unix(),
		InitialBalance:          entry.InitialBalance.BigInt(),
		SharesDst:               entry.SharesDst.BigInt(),
		UnbondingId:             entry.UnbondingId,
		UnbondingOnHoldRefCount: entry.UnbondingOnHoldRefCount,
	}
}

func convertRedelegation(redelegation stakingtypes.Redelegation) (Redelegation, error) {
	entries := make([]RedelegationEntry, len(redelegation.Entries))
	for i, v := range redelegation.Entries {
		entries[i] = convertRedelegationEntry(v)
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
	if err != nil {
		return Redelegation{}, err
	}
	validatorSrcAddress, err := sdk.ValAddressFromBech32(redelegation.ValidatorSrcAddress)
	if err != nil {
		return Redelegation{}, err
	}
	validatorDstAddress, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
	if err != nil {
		return Redelegation{}, err
	}
	return Redelegation{
		DelegatorAddress:    common.BytesToAddress(delegatorAddress.Bytes()),
		ValidatorSrcAddress: common.BytesToAddress(validatorSrcAddress.Bytes()),
		ValidatorDstAddress: common.BytesToAddress(validatorDstAddress.Bytes()),
		Entries:             entries,
	}, nil
}

func convertRedelegationEntryResponse(response stakingtypes.RedelegationEntryResponse) RedelegationEntryResponse {
	return RedelegationEntryResponse{
		RedelegationEntry: convertRedelegationEntry(response.RedelegationEntry),
		Balance:           response.Balance.BigInt(),
	}
}

func convertRedelegationResponse(response stakingtypes.RedelegationResponse) (RedelegationResponse, error) {
	entries := make([]RedelegationEntryResponse, len(response.Entries))
	for i, v := range response.Entries {
		entries[i] = convertRedelegationEntryResponse(v)
	}
	redelegation, err := convertRedelegation(response.Redelegation)
	if err != nil {
		return RedelegationResponse{}, nil
	}
	return RedelegationResponse{
		Redelegation: redelegation,
		Entries:      entries,
	}, nil
}

func convertParams(params stakingtypes.Params) Params {
	return Params{
		UnbondingTime:     int64(params.UnbondingTime.Seconds()),
		MaxValidators:     params.MaxValidators,
		MaxEntries:        params.MaxEntries,
		HistoricalEntries: params.HistoricalEntries,
		BondDenom:         params.BondDenom,
		MinCommissionRate: params.MinCommissionRate.BigInt(),
	}
}

func NewMsgCreateValidator(args []interface{}, sender common.Address, denom string) (*stakingtypes.MsgCreateValidator, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 5, len(args))
	}
	description := args[0].(Description)
	commission := args[1].(CommissionRates)
	minSelfDelegation := args[2].(*big.Int)

	pkstr := args[3].(string)
	bz, err := base64.StdEncoding.DecodeString(pkstr)
	if err != nil {
		return nil, err
	}
	var pk cryptotypes.PubKey
	if len(bz) == ed25519.PubKeySize {
		pk = &ed25519.PubKey{Key: bz}
	} else {
		return nil, errors.New(ErrPubKeyInvalidLength)
	}
	pkAny, err := codectypes.NewAnyWithValue(pk)
	if err != nil {
		return nil, err
	}

	value := args[4].(*big.Int)
	msg := &stakingtypes.MsgCreateValidator{
		Description:       convertStakingDescription(description),
		Commission:        convertStakingCommissionRates(commission),
		MinSelfDelegation: math.NewIntFromBigInt(minSelfDelegation),
		DelegatorAddress:  sdk.AccAddress(sender.Bytes()).String(),
		ValidatorAddress:  sdk.ValAddress(sender.Bytes()).String(),
		Pubkey:            pkAny,
		Value:             sdk.Coin{Denom: denom, Amount: math.NewIntFromBigInt(value)},
	}
	return msg, msg.ValidateBasic()
}

func NewMsgEditValidator(args []interface{}, sender common.Address) (*stakingtypes.MsgEditValidator, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 3, len(args))
	}
	description := args[0].(Description)

	commissionRateNullable := args[1].(NullableUint)
	var commissionRate *sdk.Dec
	if !commissionRateNullable.IsNull {
		value := precompiles_common.BigIntToLegacyDec(commissionRateNullable.Value)
		commissionRate = &value
	}

	minSelfDelegationNullable := args[2].(NullableUint)
	var minSelfDelegation *sdk.Int
	if !minSelfDelegationNullable.IsNull {
		value := math.NewIntFromBigInt(minSelfDelegationNullable.Value)
		minSelfDelegation = &value
	}

	msg := &stakingtypes.MsgEditValidator{
		Description:       convertStakingDescription(description),
		CommissionRate:    commissionRate,
		ValidatorAddress:  sdk.ValAddress(sender.Bytes()).String(),
		MinSelfDelegation: minSelfDelegation,
	}
	return msg, msg.ValidateBasic()
}

func NewMsgDelegate(args []interface{}, sender common.Address, denom string) (*stakingtypes.MsgDelegate, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	validatorAddress := args[0].(common.Address)
	amount := args[1].(*big.Int)

	msg := &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(sender.Bytes()).String(),
		ValidatorAddress: sdk.ValAddress(validatorAddress.Bytes()).String(),
		Amount:           sdk.Coin{Denom: denom, Amount: math.NewIntFromBigInt(amount)},
	}
	return msg, msg.ValidateBasic()
}

func NewMsgBeginRedelegate(args []interface{}, sender common.Address, denom string) (*stakingtypes.MsgBeginRedelegate, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 3, len(args))
	}
	validatorSrcAddress := args[0].(common.Address)
	validatorDstAddress := args[1].(common.Address)
	amount := args[2].(*big.Int)

	msg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    sdk.AccAddress(sender.Bytes()).String(),
		ValidatorSrcAddress: sdk.ValAddress(validatorSrcAddress.Bytes()).String(),
		ValidatorDstAddress: sdk.ValAddress(validatorDstAddress.Bytes()).String(),
		Amount:              sdk.Coin{Denom: denom, Amount: math.NewIntFromBigInt(amount)},
	}
	return msg, msg.ValidateBasic()
}

func NewMsgUndelegate(args []interface{}, sender common.Address, denom string) (*stakingtypes.MsgUndelegate, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	validatorAddress := args[0].(common.Address)
	amount := args[1].(*big.Int)

	msg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: sdk.AccAddress(sender.Bytes()).String(),
		ValidatorAddress: sdk.ValAddress(validatorAddress.Bytes()).String(),
		Amount:           sdk.Coin{Denom: denom, Amount: math.NewIntFromBigInt(amount)},
	}
	return msg, msg.ValidateBasic()
}

func NewMsgCancelUnbondingDelegation(args []interface{}, sender common.Address, denom string) (*stakingtypes.MsgCancelUnbondingDelegation, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 3, len(args))
	}
	validatorAddress := args[0].(common.Address)
	amount := args[1].(*big.Int)
	creationHeight := args[2].(*big.Int)

	msg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: sdk.AccAddress(sender.Bytes()).String(),
		ValidatorAddress: sdk.ValAddress(validatorAddress.Bytes()).String(),
		Amount:           sdk.Coin{Denom: denom, Amount: math.NewIntFromBigInt(amount)},
		CreationHeight:   creationHeight.Int64(),
	}
	return msg, msg.ValidateBasic()
}

func NewQueryValidatorsRequest(args []interface{}) (*stakingtypes.QueryValidatorsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	status := args[0].(string)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryValidatorsRequest{
		Status:     status,
		Pagination: convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryValidatorRequest(args []interface{}) (*stakingtypes.QueryValidatorRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 1, len(args))
	}
	validatorAddress := args[0].(common.Address)

	return &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: sdk.ValAddress(validatorAddress.Bytes()).String(),
	}, nil
}

func NewQueryValidatorDelegationsRequest(args []interface{}) (*stakingtypes.QueryValidatorDelegationsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	validatorAddr := args[0].(common.Address)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: sdk.ValAddress(validatorAddr.Bytes()).String(),
		Pagination:    convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryValidatorUnbondingDelegationsRequest(args []interface{}) (*stakingtypes.QueryValidatorUnbondingDelegationsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	validatorAddr := args[0].(common.Address)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryValidatorUnbondingDelegationsRequest{
		ValidatorAddr: sdk.ValAddress(validatorAddr.Bytes()).String(),
		Pagination:    convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryDelegationRequest(args []interface{}) (*stakingtypes.QueryDelegationRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	validatorAddr := args[1].(common.Address)

	return &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		ValidatorAddr: sdk.ValAddress(validatorAddr.Bytes()).String(),
	}, nil
}

func NewQueryUnbondingDelegationRequest(args []interface{}) (*stakingtypes.QueryUnbondingDelegationRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	validatorAddr := args[1].(common.Address)

	return &stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		ValidatorAddr: sdk.ValAddress(validatorAddr.Bytes()).String(),
	}, nil
}

func NewQueryDelegatorDelegationsRequest(args []interface{}) (*stakingtypes.QueryDelegatorDelegationsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		Pagination:    convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryDelegatorUnbondingDelegationsRequest(args []interface{}) (*stakingtypes.QueryDelegatorUnbondingDelegationsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		Pagination:    convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryRedelegationsRequest(args []interface{}) (*stakingtypes.QueryRedelegationsRequest, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 4, len(args))
	}
	delegatorAddress := args[0].(common.Address)
	validatorSrcAddress := args[1].(common.Address)
	validatorDstAddress := args[2].(common.Address)
	pagination := args[3].(PageRequest)

	return &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    sdk.AccAddress(delegatorAddress.Bytes()).String(),
		SrcValidatorAddr: sdk.ValAddress(validatorSrcAddress.Bytes()).String(),
		DstValidatorAddr: sdk.ValAddress(validatorDstAddress.Bytes()).String(),
		Pagination:       convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryDelegatorValidatorsRequest(args []interface{}) (*stakingtypes.QueryDelegatorValidatorsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	pagination := args[1].(PageRequest)

	return &stakingtypes.QueryDelegatorValidatorsRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		Pagination:    convertQueryPageRequest(pagination),
	}, nil
}

func NewQueryDelegatorValidatorRequest(args []interface{}) (*stakingtypes.QueryDelegatorValidatorRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 2, len(args))
	}
	delegatorAddr := args[0].(common.Address)
	validatorAddr := args[1].(common.Address)

	return &stakingtypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: sdk.AccAddress(delegatorAddr.Bytes()).String(),
		ValidatorAddr: sdk.ValAddress(validatorAddr.Bytes()).String(),
	}, nil
}

func NewQueryPoolRequest(args []interface{}) (*stakingtypes.QueryPoolRequest, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 0, len(args))
	}

	return &stakingtypes.QueryPoolRequest{}, nil
}

func NewQueryParamsRequest(args []interface{}) (*stakingtypes.QueryParamsRequest, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(precompiles_common.ErrInvalidNumberOfArgs, 0, len(args))
	}

	return &stakingtypes.QueryParamsRequest{}, nil
}

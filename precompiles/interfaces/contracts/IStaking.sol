// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.0;

/**
 * @dev Description defines a validator description
 */
struct Description {
    string moniker;
    string identity;
    string website;
    string securityContact;
    string details;
}

/**
 * @dev CommissionRates defines the initial commission rates to be used for creating
 * a validator.
 */
struct CommissionRates {
    uint rate; // 18 decimals
    uint maxRate; // 18 decimals
    uint maxChangeRate; // 18 decimals
}

/**
 * @dev Commission defines the commission parameters.
 */
struct Commission {
    CommissionRates commissionRates;
    uint updateTime;
}

/**
 * @dev Validator defines a validator.
 */
struct Validator {
    address operatorAddress;
    string consensusPubkey;
    bool jailed;
    BondStatus status;
    uint tokens;
    uint delegatorShares; // 18 decimals
    Description description;
    int64 unbondingHeight;
    int64 unbondingTime;
    Commission commission;
    uint minSelfDelegation;
    int64 unbondingOnHoldRefCount;
    uint64[] unbondingIds;
}

/**
 * @dev Delegation represents the bond with tokens held by an account.
 */
struct Delegation {
    address delegatorAddress;
    address validatorAddress;
    uint shares; // 18 decimals
}

/**
 * @dev RedelegationResponse is equivalent to a Redelegation except that its entries
 * contain a balance in addition to shares which is more suitable for client
 * responses.
 */
struct DelegationResponse {
    Delegation delegation;
    uint balance;
}

/**
 * @dev UnbondingDelegationEntry defines an unbonding object with relevant metadata.
 */
struct UnbondingDelegationEntry {
    int64 creationHeight;
    int64 completionTime;
    uint initialBalance;
    uint balance;
    uint64 unbondingId;
    int64 unbondingOnHoldRefCount;
}

/**
 * @dev UnbondingDelegation stores all of a single delegator's unbonding bonds
 * for a single validator in an time-ordered list.
 */
struct UnbondingDelegation {
    address delegatorAddress;
    address validatorAddress;
    UnbondingDelegationEntry[] entries;
}

/**
 * @dev RedelegationResponse is equivalent to a Redelegation except that its entries
 * contain a balance in addition to shares which is more suitable for client
 * responses.
 */
struct RedelegationResponse {
    Redelegation redelegation;
    RedelegationEntryResponse[] entries;
}

/**
 * @dev Redelegation contains the list of a particular delegator's redelegating bonds
 * from a particular source validator to a particular destination validator.
 */
struct Redelegation {
    address delegatorAddress;
    address validatorSrcAddress;
    address validatorDstAddress;
    RedelegationEntry[] entries;
}

/**
 * @dev RedelegationEntry defines a redelegation object with relevant metadata.
 */
struct RedelegationEntry {
    int64 creationHeight;
    int64 completionTime;
    uint initialBalance;
    uint sharesDst; // 18 decimals
    uint64 unbondingId;
    int64 unbondingOnHoldRefCount;
}

/**
 * @dev RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
 * contains a balance in addition to shares which is more suitable for client
 * responses.
 */
struct RedelegationEntryResponse {
    RedelegationEntry redelegationEntry;
    uint balance;
}

/**
 * @dev Params defines the parameters for the x/staking module.
 */
struct Params {
    int64 unbondingTime;
    uint32 maxValidators;
    uint32 maxEntries;
    uint32 historicalEntries;
    string bondDenom;
    uint minCommissionRate; // 18 decimals
}

/**
 * @dev BondStatus is the status of a validator.
 */
enum BondStatus {
    Unspecified,
    Unbonded,
    Unbonding,
    Bonded
}

struct NullableUint {
    bool isNull;
    uint value;
}

struct PageRequest {
    bytes key;
    uint64 offset;
    uint64 limit;
    bool countTotal;
    bool reverse;
}

struct PageResponse {
    bytes nextKey;
    uint64 total;
}

interface IStaking {
    /*=== cosmos tx ===*/

    /**
     * @dev CreateValidator defines a method for creating a new validator for tx sender.
     * cosmos grpc: rpc CreateValidator(MsgCreateValidator) returns (MsgCreateValidatorResponse);
     */
    function createValidator(
        Description memory description,
        CommissionRates memory commission,
        uint minSelfDelegation,
        string memory pubkey, // 0gchaind tendermint show-validator
        uint value
    ) external;

    /**
     * @dev EditValidator defines a method for editing an existing validator (tx sender).
     * cosmos grpc: rpc EditValidator(MsgEditValidator) returns (MsgEditValidatorResponse);
     */
    function editValidator(
        Description memory description,
        NullableUint memory commissionRate,
        NullableUint memory minSelfDelegation
    ) external;

    /**
     * @dev Delegate defines a method for performing a delegation of coins from a delegator to a validator.abi
     * The delegator is tx sender.
     * cosmos grpc: rpc Delegate(MsgDelegate) returns (MsgDelegateResponse);
     */
    function delegate(
        address validatorAddress,
        uint amount // in bond denom
    ) external;

    /**
     * @dev BeginRedelegate defines a method for performing a redelegationA
     * of coins from a delegator and source validator to a destination validator.
     * The delegator is tx sender.
     * cosmos grpc: rpc BeginRedelegate(MsgBeginRedelegate) returns (MsgBeginRedelegateResponse);
     */
    function beginRedelegate(
        address validatorSrcAddress,
        address validatorDstAddress,
        uint amount // in bond denom
    ) external returns (uint completionTime);

    /**
     * @dev Undelegate defines a method for performing an undelegation from a
     * delegate and a validator.
     * The delegator is tx sender.
     * cosmos grpc: rpc Undelegate(MsgUndelegate) returns (MsgUndelegateResponse);
     */
    function undelegate(
        address validatorAddress,
        uint amount // in bond denom
    ) external returns (uint completionTime);

    /**
     * @dev CancelUnbondingDelegation defines a method for performing canceling the unbonding delegation
     * and delegate back to previous validator.
     * The delegator is tx sender.
     * Since: cosmos-sdk 0.46
     * cosmos grpc: rpc CancelUnbondingDelegation(MsgCancelUnbondingDelegation) returns (MsgCancelUnbondingDelegationResponse);
     */
    function cancelUnbondingDelegation(
        address validatorAddress,
        uint amount, // in bond denom
        uint creationHeight
    ) external;

    /**
     * @dev UpdateParams defines an operation for updating the x/staking module parameters.
     * Since: cosmos-sdk 0.47
     * grpc: rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
     */
    // Skipped. This function is controlled by governance module.

    /*=== cosmos query ===*/

    /**
     * @dev Validators queries all validators that match the given status.
     * cosmos grpc: rpc Validators(QueryValidatorsRequest) returns (QueryValidatorsResponse);
     */
    function validators(
        string memory status,
        PageRequest memory pagination
    )
        external
        view
        returns (
            Validator[] memory validators,
            PageResponse memory paginationResult
        );

    /**
     * @dev Validator queries validator info for given validator address.
     * cosmos grpc: rpc Validator(QueryValidatorRequest) returns (QueryValidatorResponse);
     */
    function validator(
        address validatorAddress
    ) external view returns (Validator memory validator);

    /**
     * @dev ValidatorDelegations queries delegate info for given validator.
     * cosmos grpc: rpc ValidatorDelegations(QueryValidatorDelegationsRequest) returns (QueryValidatorDelegationsResponse);
     */
    function validatorDelegations(
        address validatorAddr,
        PageRequest memory pagination
    )
        external
        view
        returns (
            DelegationResponse[] memory delegationResponses,
            PageResponse memory paginationResult
        );

    /**
     * @dev ValidatorUnbondingDelegations queries unbonding delegations of a validator.
     * cosmos grpc: rpc ValidatorUnbondingDelegations(QueryValidatorUnbondingDelegationsRequest) returns (QueryValidatorUnbondingDelegationsResponse);
     */
    //
    function validatorUnbondingDelegations(
        address validatorAddr,
        PageRequest memory pagination
    )
        external
        view
        returns (
            UnbondingDelegation[] memory unbondingResponses,
            PageResponse memory paginationResult
        );

    /**
     * @dev Delegation queries delegate info for given validator delegator pair.
     * cosmos grpc: rpc Delegation(QueryDelegationRequest) returns (QueryDelegationResponse);
     */
    function delegation(
        address delegatorAddr,
        address validatorAddr
    ) external view returns (Delegation memory delegation, uint balance);

    /**
     * @dev UnbondingDelegation queries unbonding info for given validator delegator pair.
     * cosmos grpc: rpc UnbondingDelegation(QueryUnbondingDelegationRequest) returns (QueryUnbondingDelegationResponse);
     */
    function unbondingDelegation(
        address delegatorAddr,
        address validatorAddr
    ) external view returns (UnbondingDelegation memory unbond);

    /**
     * @dev DelegatorDelegations queries all delegations of a given delegator address.
     *
     * cosmos grpc: rpc DelegatorDelegations(QueryDelegatorDelegationsRequest) returns (QueryDelegatorDelegationsResponse);
     */
    function delegatorDelegations(
        address delegatorAddr,
        PageRequest memory pagination
    )
        external
        view
        returns (
            DelegationResponse[] memory delegationResponses,
            PageResponse memory paginationResult
        );

    /**
     * @dev DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address.
     * cosmos grpc: rpc DelegatorUnbondingDelegations(QueryDelegatorUnbondingDelegationsRequest)
     */
    function delegatorUnbondingDelegations(
        address delegatorAddr,
        PageRequest memory pagination
    )
        external
        view
        returns (
            UnbondingDelegation[] memory unbondingResponses,
            PageResponse memory paginationResult
        );

    /**
     * @dev Redelegations queries redelegations of given address.
     *
     * grpc: rpc Redelegations(QueryRedelegationsRequest) returns (QueryRedelegationsResponse);
     */
    function redelegations(
        address delegatorAddress,
        address srcValidatorAddress,
        address dstValidatorAddress,
        PageRequest calldata pageRequest
    )
        external
        view
        returns (
            RedelegationResponse[] calldata redelegationResponses,
            PageResponse calldata paginationResult
        );

    /**
     * @dev DelegatorValidators queries all validators info for given delegator address.
     * cosmos grpc: rpc DelegatorValidators(QueryDelegatorValidatorsRequest) returns (QueryDelegatorValidatorsResponse);
     */
    function delegatorValidators(
        address delegatorAddr,
        PageRequest memory pagination
    )
        external
        view
        returns (
            Validator[] memory validators,
            PageResponse memory paginationResult
        );

    /**
     * @dev DelegatorValidator queries validator info for given delegator validator pair.
     * cosmos grpc: rpc DelegatorValidator(QueryDelegatorValidatorRequest) returns (QueryDelegatorValidatorResponse);
     */
    function delegatorValidator(
        address delegatorAddr,
        address validatorAddr
    ) external view returns (Validator memory validator);

    /**
     * @dev Pool queries the pool info.
     * cosmos grpc: rpc Pool(QueryPoolRequest) returns (QueryPoolResponse);
     */
    function pool()
        external
        view
        returns (uint notBondedTokens, uint bondedTokens);

    /**
     * @dev Parameters queries the staking parameters.
     * cosmos grpc: rpc Params(QueryParamsRequest) returns (QueryParamsResponse);
     */
    function params() external view returns (Params memory params);
}

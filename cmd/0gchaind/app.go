package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Kava-Labs/opendb"
	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	ethermintflags "github.com/evmos/ethermint/server/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/0glabs/0g-chain/app"
	"github.com/0glabs/0g-chain/app/params"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

const (
	flagMempoolEnableAuth    = "mempool.enable-authentication"
	flagMempoolAuthAddresses = "mempool.authorized-addresses"
	flagSkipLoadLatest       = "skip-load-latest"
)

var accountNonceOp app.AccountNonceOp

// appCreator holds functions used by the sdk server to control the 0g-chain app.
// The methods implement types in cosmos-sdk/server/types
type appCreator struct {
	encodingConfig params.EncodingConfig
}

// newApp loads config from AppOptions and returns a new app.
func (ac appCreator) newApp(
	logger log.Logger,
	db cometbftdb.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	var cache sdk.MultiStorePersistentCache
	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	snapshotDir := filepath.Join(homeDir, "data", "snapshots") // TODO can these directory names be imported from somewhere?
	snapshotDB, err := opendb.OpenDB(appOpts, snapshotDir, "metadata", server.GetAppDBBackend(appOpts))
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	mempoolEnableAuth := cast.ToBool(appOpts.Get(flagMempoolEnableAuth))
	mempoolAuthAddresses, err := accAddressesFromBech32(
		cast.ToStringSlice(appOpts.Get(flagMempoolAuthAddresses))...,
	)
	if err != nil {
		panic(fmt.Sprintf("could not get authorized address from config: %v", err))
	}

	iavlDisableFastNode := appOpts.Get(server.FlagDisableIAVLFastNode)
	if iavlDisableFastNode == nil {
		iavlDisableFastNode = true
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	// Setup chainId
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if len(chainID) == 0 {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}
		chainID = appGenesis.ChainID
	}

	skipLoadLatest := false
	if appOpts.Get(flagSkipLoadLatest) != nil {
		skipLoadLatest = cast.ToBool(appOpts.Get(flagSkipLoadLatest))
	}

	bApp := app.NewBaseApp(logger, db, ac.encodingConfig,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(strings.Replace(cast.ToString(appOpts.Get(server.FlagMinGasPrices)), ";", ",", -1)),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))), // TODO what is this?
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(iavlDisableFastNode)),
		baseapp.SetIAVLLazyLoading(cast.ToBool(appOpts.Get(server.FlagIAVLLazyLoading))),
		baseapp.SetChainID(chainID),
		baseapp.SetTxInfoExtracter(extractTxInfo),
	)

	mempool := app.NewPriorityMempool(
		app.PriorityNonceWithMaxTx(fixMempoolSize(appOpts)),
		app.PriorityNonceWithTxReplacedCallback(func(ctx context.Context, oldTx, newTx *app.TxInfo) {
			if oldTx.Sender != newTx.Sender {
				sdkContext := sdk.UnwrapSDKContext(ctx)
				if accountNonceOp != nil {
					nonce := accountNonceOp.GetAccountNonce(sdkContext, oldTx.Sender)
					if nonce > 0 {
						accountNonceOp.SetAccountNonce(sdkContext, oldTx.Sender, nonce-1)
						sdkContext.Logger().Debug("rewind the nonce of the account", "account", oldTx.Sender, "from", nonce, "to", nonce-1)
					} else {
						sdkContext.Logger().Info("First meeting account", "account", oldTx.Sender)
					}
				}
			} else {
				sdkContext := sdk.UnwrapSDKContext(ctx)
				sdkContext.Logger().Info("tx replace", "account", oldTx.Sender, "nonce", oldTx.Nonce)
			}
			bApp.RegisterMempoolTxReplacedEvent(ctx, oldTx.Tx, newTx.Tx)
		}),
	)
	bApp.SetMempool(mempool)

	bApp.SetTxEncoder(ac.encodingConfig.TxConfig.TxEncoder())
	abciProposalHandler := app.NewDefaultProposalHandler(mempool, bApp)
	bApp.SetPrepareProposal(abciProposalHandler.PrepareProposalHandler())

	newApp := app.NewApp(
		homeDir, traceStore, ac.encodingConfig,
		app.Options{
			SkipLoadLatest:        skipLoadLatest,
			SkipUpgradeHeights:    skipUpgradeHeights,
			SkipGenesisInvariants: cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants)),
			InvariantCheckPeriod:  cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
			MempoolEnableAuth:     mempoolEnableAuth,
			MempoolAuthAddresses:  mempoolAuthAddresses,
			EVMTrace:              cast.ToString(appOpts.Get(ethermintflags.EVMTracer)),
			EVMMaxGasWanted:       cast.ToUint64(appOpts.Get(ethermintflags.EVMMaxTxGasWanted)),
		},
		bApp,
	)

	accountNonceOp = app.NewAccountNonceOp(newApp)

	return newApp
}

// appExport writes out an app's state to json.
func (ac appCreator) appExport(
	logger log.Logger,
	db cometbftdb.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	options := app.DefaultOptions
	options.SkipLoadLatest = true
	options.InvariantCheckPeriod = cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	var tempApp *app.App
	if height != -1 {
		bApp := app.NewBaseApp(logger, db, ac.encodingConfig)
		tempApp = app.NewApp(homePath, traceStore, ac.encodingConfig, options, bApp)

		if err := tempApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		bApp := app.NewBaseApp(logger, db, ac.encodingConfig)
		tempApp = app.NewApp(homePath, traceStore, ac.encodingConfig, options, bApp)
	}
	return tempApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

// addStartCmdFlags adds flags to the server start command.
func (ac appCreator) addStartCmdFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// accAddressesFromBech32 converts a slice of bech32 encoded addresses into a slice of address types.
func accAddressesFromBech32(addresses ...string) ([]sdk.AccAddress, error) {
	var decodedAddresses []sdk.AccAddress
	for _, s := range addresses {
		a, err := sdk.AccAddressFromBech32(s)
		if err != nil {
			return nil, err
		}
		decodedAddresses = append(decodedAddresses, a)
	}
	return decodedAddresses, nil
}

var ErrMustHaveSigner error = errors.New("tx must have at least one signer")

func extractTxInfo(ctx sdk.Context, tx sdk.Tx) (*sdk.TxInfo, error) {
	sigs, err := tx.(signing.SigVerifiableTx).GetSignaturesV2()
	if err != nil {
		return nil, err
	}

	var sender string
	var nonce uint64
	var gasPrice uint64
	var gasLimit uint64
	var txType int32

	if len(sigs) == 0 {
		txType = 1
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return nil, ErrMustHaveSigner
		}
		msgEthTx, ok := msgs[0].(*evmtypes.MsgEthereumTx)
		if !ok {
			return nil, ErrMustHaveSigner
		}
		ethTx := msgEthTx.AsTransaction()
		signer := gethtypes.NewEIP2930Signer(ethTx.ChainId())
		ethSender, err := signer.Sender(ethTx)
		if err != nil {
			return nil, ErrMustHaveSigner
		}
		sender = sdk.AccAddress(ethSender.Bytes()).String()
		nonce = ethTx.Nonce()
		gasPrice = ethTx.GasPrice().Uint64()
		gasLimit = ethTx.Gas()
	} else {
		sig := sigs[0]
		sender = sdk.AccAddress(sig.PubKey.Address()).String()
		nonce = sig.Sequence
	}

	return &sdk.TxInfo{
		SignerAddress: sender,
		Nonce:         nonce,
		GasLimit:      gasLimit,
		GasPrice:      gasPrice,
		TxType:        txType,
	}, nil
}

func fixMempoolSize(appOpts servertypes.AppOptions) int {
	val1 := appOpts.Get("mempool.size")
	val2 := appOpts.Get(server.FlagMempoolMaxTxs)

	if val1 != nil && val2 != nil {
		size1 := cast.ToInt(val1)
		size2 := cast.ToInt(val2)
		if size1 != size2 {
			panic("the value of mempool.size and mempool.max-txs are different")
		}
		return size1
	} else if val1 == nil && val2 == nil {
		panic("not found mempool size in config")
	} else if val1 == nil {
		return cast.ToInt(val2)
	} else { //if val2 == nil {
		return cast.ToInt(val1)
	}
}

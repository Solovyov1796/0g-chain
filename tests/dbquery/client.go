package main

import (
	"encoding/hex"
	"strings"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dbName = "dbQuery"
)

var (
	appDB       dbm.DB
	blockDB     dbm.DB
	stateDB     dbm.DB
	indexDB     dbm.DB
	encodingCfg EncodingConfig
)

func InitDB(backend dbm.BackendType, dir string) error {
	encodingCfg = MakeEncodingConfig()
	var err error
	appDB, err = dbm.NewDB("application", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create application db")
	}
	blockDB, err = dbm.NewDB("blockstore", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create blockstore db")
	}
	stateDB, err = dbm.NewDB("state", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create state db")
	}
	indexDB, err = dbm.NewDB("tx_index", backend, dir)
	if err != nil {
		return errors.Wrapf(err, "failed to create tx_index db")
	}
	return nil
}

func GetTx(hash string) (sdk.Tx, error) {
	if len(hash) == 0 {
		return nil, errors.New("tx hash cannot be empty")
	}

	txByteSlice, err := queryTx(hash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Errorf(codes.NotFound, "tx not found: %s", hash)
		}

		return nil, err
	}

	// input tx byte slice
	tx, err := encodingCfg.TxConfig.TxDecoder()(txByteSlice)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func queryTx(hashHexStr string) ([]byte, error) {
	r, err := loadTxFromIndex(hashHexStr)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, errors.New("tx not found")
	}

	return r.Tx, nil
}

func loadTxFromIndex(hashHexStr string) (*abci.TxResult, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	rawBytes, err := indexDB.Get(hash)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get tx from indexDB")
	}

	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(abci.TxResult)
	err = proto.Unmarshal(rawBytes, txResult)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal txResult")
	}

	jsonStr, _ := AbciTxResultToJSONString(txResult)

	println("abci.TxResult: ", jsonStr)

	return txResult, nil
}

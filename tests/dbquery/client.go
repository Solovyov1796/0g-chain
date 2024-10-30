package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/syndtr/goleveldb/leveldb"
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

	printLastTx(dir)
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

func printLastTx(dir string) {
	// 打开 blockstore 数据库
	db, err := leveldb.OpenFile(path.Join(dir, "blockstore.db"), nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 遍历区块数据库查找最后一个区块
	iter := db.NewIterator(nil, nil)
	var lastBlockKey, lastBlockValue []byte
	for iter.Next() {
		lastBlockKey = iter.Key()
		lastBlockValue = iter.Value()
	}
	iter.Release()

	// 检查迭代器错误
	if err := iter.Error(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}

	// 假设交易数据存储在最后一个区块数据中，解析并输出
	var blockData map[string]interface{}
	if err := json.Unmarshal(lastBlockValue, &blockData); err != nil {
		log.Fatalf("Failed to unmarshal block data: %v", err)
	}

	// 输出区块信息，获取交易数据
	fmt.Printf("Last block key: %s\n", lastBlockKey)
	fmt.Printf("Last block data: %+v\n", blockData)

	// 获取并解析交易数据
	// 假设交易数据位于 blockData["txs"]
	txs, ok := blockData["txs"].([]interface{})
	if !ok {
		log.Fatalf("Transaction data not found in last block.")
	}
	lastTx := txs[len(txs)-1]

	// 输出最后一个交易
	fmt.Printf("Last transaction: %+v\n", lastTx)

}

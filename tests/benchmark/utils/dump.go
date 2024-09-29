package utils

import (
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionOutput struct {
	Nonce    uint64   `json:"nonce"`
	To       string   `json:"to"`
	Value    *big.Int `json:"value"`
	GasLimit uint64   `json:"gas_limit"`
	GasPrice *big.Int `json:"gas_price"`
	Data     []byte   `json:"data"`
}

func DumpTx(tx *types.Transaction) string {
	txOutput := TransactionOutput{
		Nonce:    tx.Nonce(),
		To:       tx.To().Hex(),
		Value:    tx.Value(),
		GasLimit: tx.Gas(),
		GasPrice: tx.GasPrice(),
		Data:     tx.Data(),
	}

	jsonData, err := json.MarshalIndent(txOutput, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal transaction to JSON: %v", err)
	}

	return string(jsonData)
}

package main

import (
	"encoding/json"
	"fmt"

	abcitypes "github.com/cometbft/cometbft/abci/types"
)

func AbciTxResultToJSONString(txResult *abcitypes.TxResult) (string, error) {
	txResultJSON, err := json.Marshal(txResult)
	if err != nil {
		return "", fmt.Errorf("failed to marshal TxResult to JSON: %w", err)
	}

	return string(txResultJSON), nil
}

package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"cosmossdk.io/errors"
	"github.com/0glabs/0g-chain/tests/benchmark/producer"
	"github.com/0glabs/0g-chain/tests/benchmark/utils"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

var senders []*Sender

func main() {
	rpcUrl, _ := getParameters()

	faucetPk, err := getAccountPrivateKey(evmFaucetMnemonic)
	if err != nil {
		log.Fatalf("Failed to get the faucet private key: %v", err)
	}
	faucetPkStr := hex.EncodeToString(faucetPk.D.Bytes())

	senders = make([]*Sender, 0, senderCount)
	for i := range senderCount {
		client, err := ethclient.Dial(rpcUrl)
		if err != nil {
			log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		}

		generator, err := producer.NewTransferGenerator(uint32(basePrefix+i), accountCount, faucetPkStr, client)
		// generator, err := producer.NewErc20Generator(accountCount, "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306", client)
		if err != nil {
			log.Fatalf("[%d]Failed to create the generator: %v", i, err.Error())
		}

		thisSender := &Sender{
			Index:     i,
			Generator: generator,
			Client:    client,
		}
		senders = append(senders, thisSender)

		err = generator.WarmUp()
		if err != nil {
			log.Fatalf("[%d]Failed to warm up the generator: %v", i, err.Error())
		}

		tx := make([](*types.Transaction), 0, accountCount)
		thisSender.SendCh = generator.GenerateTransfer()

		for t := range thisSender.SendCh {
			tx = append(tx, t)
			if len(tx) == accountCount {
				break
			}
		}
		sendWorkload(client, tx)
		time.Sleep(10 * time.Second)
	}

	// time.Sleep(30 * time.Second)

	println("ready!!!", time.Now().Format("2006-01-02 15:04:05"))

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	startBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get the start block: %v", err)
	}

	swg := utils.NewSizedWaitGroup(senderCount)
	for i := range senders {
		swg.Add()
		go func(s *Sender) {
			defer swg.Done()
			s.Send(txCount)
		}(senders[i])
	}

	swg.Wait()

	println("done!!!", time.Now().Format("2006-01-02 15:04:05"))

	endBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get the end block: %v", err)
	}

	startHeight := startBlock.Number()
	startTime := startBlock.Time()

	endHeight := endBlock.Number()
	endTime := endBlock.Time()

	if endTime <= startTime {
		log.Fatal("More transactions are needed")
	}

	totalTxCount := 0

	// Iterate over the blocks from startHeight to the endHeight
	for blockNumber := startHeight.Add(startHeight, big.NewInt(1)); blockNumber.Cmp(endHeight) <= 0; blockNumber.Add(blockNumber, big.NewInt(1)) {
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatalf("Failed to fetch block: %v", err)
		}

		totalTxCount += len(block.Transactions())
	}

	elapsedSeconds := endTime - startTime
	tps := float64(totalTxCount) / float64(elapsedSeconds)

	fmt.Printf("Total transactions counted in %d seconds is %d\n", elapsedSeconds, totalTxCount)
	fmt.Printf("The TPS of the chain is %.2f\n", tps)
}

func getParameters() (string, int) {
	// handle command line flags
	rpcUrl := flag.String("rpc-url", "http://127.0.0.1:8545", "RPC url of the chain")
	count := flag.Int("count", 10000, "The number of transactions to be sent")
	flag.Parse()

	if *count > 1000000 {
		log.Fatal("Too many transactions to be generated and sent")
	}

	return *rpcUrl, *count
}

func sendWorkload(client *ethclient.Client, workload [](*types.Transaction)) {
	for _, tx := range workload {
		err := client.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Fatal("Failed to send transactions ", tx.Hash().String(), " error: ", err.Error())
		} else {
			println("Sent transaction", tx.Hash().String())
		}
	}
}

func getAccountPrivateKey(mnemonic string) (*ecdsa.PrivateKey, error) {
	hdPath := hd.CreateHDPath(60, 0, 0)
	privKeyBytes, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to derive private key from mnemonic")
	}
	privKey := &ethsecp256k1.PrivKey{Key: privKeyBytes}
	return crypto.HexToECDSA(hex.EncodeToString(privKey.Bytes()))
}

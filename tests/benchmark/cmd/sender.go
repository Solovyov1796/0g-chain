package cmd

import (
	"context"
	"log"
	"time"

	"github.com/0glabs/0g-chain/tests/benchmark/generator"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/time/rate"
)

type Sender struct {
	Speed     int
	Generator generator.Generator
	Client    *ethclient.Client
	SendCh    <-chan *types.Transaction
}

func (s *Sender) Send() {
	ctx := context.Background()
	limiter := rate.NewLimiter(rate.Limit(s.Speed), 1)
	for t := range s.SendCh {
		limiter.Wait(context.Background())
		err := s.Client.SendTransaction(ctx, t)
		if err != nil {
			log.Fatal("Failed to send transactions ", t.Hash().String(), " error: ", err.Error())
		} else {
			println(time.Now().Format("2006-01-02 15:04:05.000000"), ">>> ", "Sent transaction", t.Hash().String())
		}
	}
}

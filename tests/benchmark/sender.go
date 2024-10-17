package main

import (
	"context"
	"log"
	"time"

	"github.com/0glabs/0g-chain/tests/benchmark/producer"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/time/rate"
)

type Sender struct {
	Index     int
	Generator producer.Generator
	Client    *ethclient.Client
	SendCh    <-chan *types.Transaction
}

func (s *Sender) Send(txCount int) {
	ctx := context.Background()
	cnt := 0
	limiter := rate.NewLimiter(rate.Every(time.Second/senderSpeed), 1)
	for t := range s.SendCh {
		limiter.Wait(ctx)
		err := s.Client.SendTransaction(ctx, t)
		if err != nil {
			log.Fatal("[", s.Index, "] Failed to send transactions ", t.Hash().String(), " error: ", err.Error())
		} else {
			println("[", s.Index, "] Sent transaction", t.Hash().String())
		}

		cnt++
		if cnt == txCount {
			break
		}
	}
}

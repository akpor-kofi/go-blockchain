package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/akpor-kofi/blockchain/components"
	"github.com/akpor-kofi/blockchain/utils"
	"github.com/go-redis/redis/v9"
	"log"
	"sync"
	"time"
)

var (
	Rdb         *redis.Client
	BlockSub    *redis.PubSub
	TransactSub *redis.PubSub
)

var wg = sync.WaitGroup{}
var ctx = context.TODO()

func init() {
	Rdb = redis.NewClient(&redis.Options{
		//Addr:        "localhost:6379",
		Addr:        "rediss://:p8850c8a466c98e2df1f8a910147e24ca286977e57aa6e41b82f2d96ff51d604c@ec2-23-20-19-160.compute-1.amazonaws.com:31470",
		Password:    "",
		DB:          0,
		ReadTimeout: 24 * time.Hour,
	})

	//chain := getCurrentChain()
	//BC.ReplaceChain(*chain)
}

func BroadcastChain() {
	res, err := utils.MarshalJSON(BC.Chain)
	if err != nil {

	}
	ProduceEvents(res, "block:created", "block")
}

func BroadcastTransaction(tr *components.Transaction) {
	res, err := utils.MarshalJSON(&tr)
	if err != nil {
		log.Fatal(err)
	}

	ProduceEvents(res, "transaction:created", "transaction")

	// TODO: delete the
}

func ClosePubsub(pubsub *redis.PubSub) {
	err := pubsub.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func ProduceEvents(chainJSON string, stream string, key string) error {
	fmt.Println("producing events to chain streams")
	err := Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: 1,
		ID:     "",
		Values: []string{key, chainJSON},
	}).Err()

	return err
}

func ConsumeEvents() {
	chainStream := "block:created"
	transactionStream := "transaction:created"

	for i, j := "1656833753963-0", "0"; ; {
		entries, err := Rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{chainStream, transactionStream, i, j},
			Count:   1,
			Block:   24 * time.Hour,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(entries[0].Stream)
		values := entries[0].Messages[0].Values

		switch entries[0].Stream {
		case chainStream:
			i = entries[0].Messages[0].ID

			blJSON := fmt.Sprintf("%v", values["block"])
			bl := new([]components.Block)
			if err := json.Unmarshal([]byte(blJSON), bl); err != nil {
				log.Fatal(err)
			}

			err := BC.ReplaceChain(*bl, false, func() {
				TransactionPool.ClearBlockchainTransactions(*bl)
			})

			if err != nil {
				fmt.Println(err)
			}

		case transactionStream:
			j = entries[0].Messages[0].ID

			trJSON := fmt.Sprintf("%v", values["transaction"])
			//fmt.Println(trJSON)
			tr := new(components.Transaction)

			if err := json.Unmarshal([]byte(trJSON), tr); err != nil {
				log.Fatal(err)
			}
			TransactionPool.SetTransactions(tr)
		}
	}
	wg.Done()

}

func getCurrentChain() *[]components.Block {
	result, err := Rdb.XRevRangeN(ctx, "block:created", "+", "-", int64(1)).Result()
	if err != nil {
		fmt.Println("probably throwing error cause the stream has not been created or have anything")
	}

	chain := new([]components.Block)

	chainJSON := fmt.Sprintf("%v", result[0].Values["block"]) // converts the block interface to string

	err = json.Unmarshal([]byte(chainJSON), chain)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(chain)
	return chain
}

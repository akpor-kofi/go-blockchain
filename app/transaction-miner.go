package app

import (
	"fmt"
	"github.com/akpor-kofi/blockchain/components"
)

type TransactionMiner struct {
	Blockchain      *components.Blockchain      `json:"blockchain"`
	TransactionPool *components.TransactionPool `json:"transactionPool"`
	Wallet          *components.Wallet          `json:"wallet"`
}

type Miner interface {
	MineTransactions()
}

func (tm *TransactionMiner) Construct(bc *components.Blockchain, tp *components.TransactionPool, w *components.Wallet) *TransactionMiner {
	tm.Blockchain = bc
	tm.TransactionPool = tp
	tm.Wallet = w

	return tm
}

func (tm *TransactionMiner) MineTransaction() error {
	validTransactions := tm.TransactionPool.ValidTransactions()
	fmt.Println(validTransactions, "valid bro")
	if len(validTransactions) == 0 {
		return fmt.Errorf("%v", "there is no valid transactions to mine")
	}

	validTransactions = append(validTransactions, components.RewardTransaction(tm.Wallet))

	tm.Blockchain.AddBlock(validTransactions)

	BroadcastChain()

	tm.TransactionPool.Clear()

	return nil
}

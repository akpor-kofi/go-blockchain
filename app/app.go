package app

import (
	"fmt"
	"github.com/akpor-kofi/blockchain/components"
)

var (
	BC              *components.Blockchain
	Wallet          *components.Wallet
	TransactionPool *components.TransactionPool
	TransactMiner   *TransactionMiner
)

func init() {
	bl := components.Block{}.Genesis()

	BC = new(components.Blockchain)
	BC.Chain = []components.Block{bl}

	Wallet = new(components.Wallet).Construct()
	fmt.Println(Wallet.PublicKey, "here")
	TransactionPool = new(components.TransactionPool).Construct()
	TransactMiner = new(TransactionMiner).Construct(BC, TransactionPool, Wallet)
	//
}

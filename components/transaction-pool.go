package components

import (
	"fmt"
	"github.com/samber/lo"
)

type TransactionPool struct {
	TransactionMap map[string]Transaction `json:"transactionMap"`
}

func (tp *TransactionPool) Clear() {
	tp.TransactionMap = make(map[string]Transaction)
}

func (tp *TransactionPool) Construct() *TransactionPool {
	tp.TransactionMap = make(map[string]Transaction)
	return tp
}

func (tp *TransactionPool) SetTransactions(t *Transaction) {
	tp.TransactionMap[t.Id] = *t
}

func (tp *TransactionPool) ExistingTransaction(inputAddress string) (*Transaction, bool) {
	t, ok := lo.Find(lo.Values(tp.TransactionMap), func(transaction Transaction) bool {
		return transaction.Address == inputAddress
	})

	return &t, ok
}

func (tp *TransactionPool) ValidTransactions() []Transaction {
	return lo.Filter(lo.Values(tp.TransactionMap), func(tr Transaction, _ int) bool {
		fmt.Println(tr)
		return ValidTransaction(&tr)
	})
}

func (tp *TransactionPool) ClearBlockchainTransactions(chain []Block) {
	for i, block := range chain {
		if i == 0 {
			continue
		}

		lo.ForEach(block.Data, func(tr Transaction, _ int) {
			delete(tp.TransactionMap, tr.Id)
		})

	}
}

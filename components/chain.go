package components

import (
	"fmt"
	"github.com/akpor-kofi/blockchain/utils"
	mapset "github.com/deckarep/golang-set"
	"github.com/samber/lo"
	"reflect"
)

type Blockchain struct {
	Chain []Block
}

type ChainValidator interface {
	IsValidChain(chain []Block) bool
}

func (bc *Blockchain) AddBlock(Data []Transaction) *Blockchain {

	bl := MineBlock(bc.Chain[len(bc.Chain)-1], Data)
	bc.Chain = append(bc.Chain, bl)
	return bc
}

func (bc *Blockchain) ValidTransactionData(chain []Block) bool {
	// enforces all the key cryptocurrency rules to protect the blockchain
	for i, block := range chain {
		if i == 0 {
			continue
		}
		// making use of sets by using maps in go
		transactionSet := mapset.NewSet()
		rewardTransactCount := 1

		for _, tr := range block.Data {
			switch {
			// reward transactions
			case tr.Input.Address == RewardInputAddress:
				rewardTransactCount = rewardTransactCount + 1

				if rewardTransactCount > 1 {
					fmt.Println("miner reward exceeds limit")
					return false
				}

				if lo.Values(tr.OutputMap)[0] != MiningReward {
					fmt.Println("miner reward amount is invalid")
					return false
				}

			case tr.Input.Address != RewardInputAddress:
				if !ValidTransaction(&tr) {
					fmt.Println("invalid transaction")
					return false
				}

				trueBalance := CalculateBalance(bc.Chain, tr.Input.Address)

				if tr.Input.Amount != trueBalance {
					fmt.Println("invalid input amount")
					return false
				}

				if transactionSet.Contains(tr) {
					return false
				} else {
					transactionSet.Add(tr)
				}
			}
		}

	}

	return true
}

func (bc *Blockchain) IsValidChain() bool {
	bl := Block{}.Genesis()

	// NB: if i was representing Data as array i would compared the json marshal version of them
	if reflect.DeepEqual(bc.Chain[0], bl) {
		fmt.Println("not valid because the first block is not the standard genesis")
		return false
	}

	for i := 1; i < len(bc.Chain); i++ {
		bl := bc.Chain[i]
		alh := bc.Chain[i-1].Hash      //actual last hash
		ld := bc.Chain[i-1].Difficulty // last Difficulty
		if alh != bl.LastHash {
			return false
		}
		h := Hash{bl.Difficulty, bl.Nonce, bl.Timestamp, bl.Data, bl.LastHash}
		vh := utils.CryptoHash(h) // validated Hash

		if bl.Hash != vh {
			return false
		}

		// check to prevent Difficulty jumps
		if ld-bl.Difficulty > 1 {
			return false
		}
	}

	return true
}

func (bc *Blockchain) ReplaceChain(chain []Block, validateTransaction bool, onSuccess func()) error {
	if len(chain) <= len(bc.Chain) {
		return fmt.Errorf("%v", "can't replace blockchain: not long enough")
	}
	if !bc.IsValidChain() {
		return fmt.Errorf("%v", "can't replace blockchain: not a valid chain")
	}

	if validateTransaction && !bc.ValidTransactionData(chain) {
		return fmt.Errorf("%v", "incoming chain has invalid data")
	}

	fmt.Println("chain replaced")
	if onSuccess != nil {
		onSuccess()
	}
	bc.Chain = chain
	return nil
}

package components

import (
	"crypto/sha256"
	"fmt"
	"github.com/akpor-kofi/blockchain/utils"
	"strings"
	"time"
)

const MineRate = 1000

type Hash struct {
	Difficulty int
	Nonce      int
	Timestamp  int64
	Data       []Transaction
	LastHash   string
}

type Block struct {
	Difficulty int           `bson:"difficulty" json:"difficulty"`
	Nonce      int           `bson:"nonce" json:"nonce"`
	Timestamp  int64         `bson:"timestamp" json:"timestamp"`
	Data       []Transaction `bson:"data" json:"data"`
	LastHash   string        `bson:"lastHash" json:"lastHash"`
	Hash       string        `bson:"hash" json:"hash"`
}

type Starter interface {
	Genesis() Block
}

type Miner interface {
	MineBlock(lastBlock Block, Data []Transaction) Block
}

func (_ Block) Genesis() Block {
	h := sha256.New()
	h.Write([]byte("genesis block"))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	w := new(Wallet).Construct()
	t := new(Transaction).Construct(w, "genesis", 20)

	return Block{
		Timestamp:  1,
		Data:       []Transaction{*t}, // had to put nil for import cycles reasons
		Hash:       hash,
		LastHash:   hash,
		Difficulty: 3,
		Nonce:      0,
	}
}

func MineBlock(lastBlock Block, Data []Transaction) Block {
	t := time.Now().UnixMilli()
	lh := lastBlock.Hash
	diff := lastBlock.Difficulty
	Nonce := 0

	var hash string

	// proof of work algorithm
	for {
		Nonce++
		t = time.Now().UnixMilli()
		diff = Block{}.AdjustDifficulty(lastBlock, t)
		h := Hash{diff, Nonce, t, Data, lh}
		hash = utils.CryptoHash(h)

		if hash[0:diff] == strings.Repeat("0", diff) {
			break
		}
	}

	return Block{
		Difficulty: diff,
		Nonce:      Nonce,
		Timestamp:  t,
		Data:       Data,
		LastHash:   lh,
		Hash:       hash,
	}
}

func (_ Block) AdjustDifficulty(oldBl Block, t int64) int {
	difference := t - oldBl.Timestamp

	if oldBl.Difficulty < 1 {
		return 1
	}

	if difference > MineRate {
		return oldBl.Difficulty - 1
	}

	return oldBl.Difficulty + 1
}

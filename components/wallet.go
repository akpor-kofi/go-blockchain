package components

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/akpor-kofi/blockchain/utils"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/samber/lo"
	"log"
)

var startingBalance = 1000
var privateKey string

type Wallet struct {
	Balance   int    `json:"balance"`
	PublicKey string `json:"publicKey"`
}

func (w *Wallet) Construct() *Wallet {
	w.Balance = startingBalance

	priKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		fmt.Println("error generating key")
		return nil
	}

	pubKey := priKey.PubKey()
	privateKey, w.PublicKey = EncodeSec(priKey, pubKey)

	return w
}

func (w *Wallet) Sign(data map[string]int) string {
	msg := utils.CryptoHash(data)

	decodeMsg, err := hex.DecodeString(msg)
	if err != nil {
		log.Fatal(err)
	}

	//privateKeyECDSA, _ := DecodeKeys(privateKey, w.PublicKey)
	pri, _ := DecodeSec(privateKey)
	privateKeyECDSA := pri.ToECDSA()

	signBytes, err := ecdsa.SignASN1(rand.Reader, privateKeyECDSA, decodeMsg)
	if err != nil {
		log.Fatal(err)
	}

	return hexutil.Encode(signBytes)
}

func (w *Wallet) CreateTransaction(r string, a int, chain []Block) (*Transaction, error) {
	if chain != nil {
		w.Balance = CalculateBalance(chain, w.PublicKey)
	}

	if a > w.Balance {
		return nil, fmt.Errorf("%s", "Amount exceeds balance")
	}

	return new(Transaction).Construct(w, r, a), nil
}

func CalculateBalance(chain []Block, address string) int {
	hasConductedTransaction := false
	outputsTotal := 0

	for i := len(chain) - 1; i > 0; i-- {
		block := chain[i]

		for _, transaction := range block.Data {
			if transaction.Input.Address == address {
				hasConductedTransaction = true
			}

			addressOutput, ok := transaction.OutputMap[address]

			if ok {
				outputsTotal = outputsTotal + addressOutput
			}
		}

		if hasConductedTransaction {
			break
		}
	}

	return lo.Ternary(hasConductedTransaction, outputsTotal, startingBalance+outputsTotal)
}

func EncodeSec(privateKey *secp256k1.PrivateKey, publicKey *secp256k1.PublicKey) (string, string) {
	privBytes := privateKey.Serialize()
	pubBytes := publicKey.SerializeCompressed()
	return hexutil.Encode(privBytes), hexutil.Encode(pubBytes)
}

func DecodeSec(privateKey string) (*secp256k1.PrivateKey, *secp256k1.PublicKey) {
	privBytes, _ := hexutil.Decode(privateKey)
	privKey := secp256k1.PrivKeyFromBytes(privBytes)
	return privKey, privKey.PubKey()
}

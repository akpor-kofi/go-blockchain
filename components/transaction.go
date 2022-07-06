package components

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/akpor-kofi/blockchain/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	randID "github.com/gofiber/fiber/v2/utils"
	"time"
)

const (
	RewardInputAddress = "*authorized-reward*"
	MiningReward       = 50
)

type Input struct {
	Timestamp int64  `json:"timestamp"`
	Amount    int    `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type Transaction struct {
	Id        string         `json:"id"`
	OutputMap map[string]int `json:"outputMap"`
	Input     `json:"input"`
}

func (t *Transaction) Construct(w *Wallet, r string, a int) *Transaction { // om & i are for reward transactions
	t.Id = randID.UUID()
	t.OutputMap = t.CreateOutputMap(w, r, a)
	t.Input = t.CreateInput(w)

	return t
}

func (t *Transaction) RewardTransactConstruct(om map[string]int, i Input) Transaction {
	t.OutputMap = om
	t.Input = i
	t.Id = randID.UUID()

	return *t
}

func (t *Transaction) CreateInput(w *Wallet) Input {

	return Input{
		Timestamp: time.Now().UnixMilli(),
		Amount:    w.Balance,
		Address:   w.PublicKey,
		Signature: w.Sign(t.OutputMap),
	}
}

func (t *Transaction) CreateOutputMap(w *Wallet, r string, a int) map[string]int {

	return map[string]int{
		r:           a,
		w.PublicKey: w.Balance - a,
	}
}

func (t *Transaction) Update(w *Wallet, r string, a int) {
	if a > t.OutputMap[w.PublicKey] {
		fmt.Println("Amount exceeds balance")
		return
	}

	if _, ok := t.OutputMap[r]; ok == false {
		t.OutputMap[r] = a
	} else {
		t.OutputMap[r] = t.OutputMap[r] + a
	}

	t.OutputMap[w.PublicKey] = t.OutputMap[w.PublicKey] - a

	t.Input = t.CreateInput(w)
}

func ValidTransaction(t *Transaction) bool {
	outputTotal := 0
	for _, v := range t.OutputMap {
		outputTotal += v
	}

	if t.Input.Amount != outputTotal {
		fmt.Println("the amount is not equal: fraudulent practice alert")
		return false
	}

	if !verifySignature(t.Signature, t.OutputMap) {
		fmt.Println("the signature cannot be verified")
		return false
	}

	return true
}

func RewardTransaction(minerWallet *Wallet) Transaction {
	input := Input{Address: RewardInputAddress}
	om := map[string]int{
		minerWallet.PublicKey: MiningReward,
	}

	t := new(Transaction).RewardTransactConstruct(om, input)

	return t

}

func verifySignature(signature string, data interface{}) bool {
	decodeMsg, _ := hex.DecodeString(utils.CryptoHash(data)) // data in bytes
	decodeSig, _ := hexutil.Decode(signature)
	//_, publicKey := DecodeKeys(privateKey, pubKey)
	_, publicKey := DecodeSec(privateKey)
	publicKeyECDSA := publicKey.ToECDSA()

	//return ecdsa.VerifyASN1(publicKey, decodeMsg, decodeSig)
	return ecdsa.VerifyASN1(publicKeyECDSA, decodeMsg, decodeSig)
}

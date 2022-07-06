package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

func MarshalJSON(c interface{}) (string, error) {
	res, err := json.Marshal(c)

	if err != nil {
		return "", err
	}

	return string(res), nil
}

func CryptoHash(o interface{}) string {
	h := sha256.New()
	// h.Write([]byte(fmt.Sprintf("%v", o)))
	marshal, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	h.Write(marshal)

	hash := hex.EncodeToString(h.Sum(nil))

	//fmt.Println(len(h.Sum(nil)))   // 32 bytes
	//fmt.Println(len(hex.DecodeString(h.Sum(nil))))   // 32 bytes
	//fmt.Println(len([]byte(hash))) // 64 bytes

	return hash
}

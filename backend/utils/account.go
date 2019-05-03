package utils

import (
	cRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
)

// GenereateKeys will generate public and private key pairs
// The public key is the account
func GenereateKeys() (string, string) {
	randomBytes := [32]byte{}
	_, err := io.ReadFull(cRand.Reader, randomBytes[:])

	if err != nil {
		fmt.Println("Failed to get randomness for the private key...")
		return "", ""
	}
	priKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("Failed to generate private key")
		return "", ""
	}

	privateKey := crypto.FromECDSA(priKey)
	address := crypto.PubkeyToAddress(priKey.PublicKey)

	return address.Hex(), hex.EncodeToString(privateKey)
}

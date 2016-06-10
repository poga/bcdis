// an account is a ESDSA public/private key pair
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"

	"github.com/tv42/base58"

	"golang.org/x/crypto/ripemd160"
)

type Account struct {
	Key *ecdsa.PrivateKey
}

func (a *Account) Public() *ecdsa.PublicKey {
	return &a.Key.PublicKey
}

func (a *Account) Address() ([]byte, error) {
	pub := ecdsaPublicKey{X: a.Public().X, Y: a.Public().Y}

	serialized, err := asn1.Marshal(pub)
	if err != nil {
		return []byte{}, err
	}

	hash := sha256.Sum256(serialized)
	bytes := ripemd160.New().Sum(hash[:])

	bigInt := new(big.Int).SetBytes(bytes)

	return base58.EncodeBig([]byte{}, bigInt), nil
}

type ecdsaPublicKey struct {
	X *big.Int
	Y *big.Int
}

func NewAccount() (*Account, error) {
	key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Account{key}, nil
}

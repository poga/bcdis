// an account is a ESDSA public/private key pair
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type Account struct {
	Key *ecdsa.PrivateKey
}

func (a *Account) Public() *ecdsa.PublicKey {
	return &a.Key.PublicKey
}

func NewAccount() (*Account, error) {
	key, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Account{key}, nil
}

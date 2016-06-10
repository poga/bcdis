package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"errors"
	"math/big"
)

type Signable interface {
	Hash() ([32]byte, error)
	SignWith(signature []byte) error
	Signature() ([]byte, error)
}

func Sign(signable Signable, account *Account) error {
	hash, err := signable.Hash()
	if err != nil {
		return err
	}

	r, s, err := ecdsa.Sign(rand.Reader, account.Key, hash[:])
	if err != nil {
		return err
	}

	sig, err := asn1.Marshal(signature{r, s})
	if err != nil {
		return err
	}
	return signable.SignWith(sig)
}

func Verify(signable Signable, account *Account) error {
	hash, err := signable.Hash()
	if err != nil {
		return err
	}

	sig, err := signable.Signature()
	if err != nil {
		return err
	}

	var ecdsaSignature signature
	_, err = asn1.Unmarshal(sig, &ecdsaSignature)
	if err != nil {
		return err
	}

	ok := ecdsa.Verify(account.Public(), hash[:], ecdsaSignature.R, ecdsaSignature.S)
	if !ok {
		return errors.New("Verification Failed")
	}

	return nil
}

// ECDSA signature
type signature struct {
	R *big.Int
	S *big.Int
}

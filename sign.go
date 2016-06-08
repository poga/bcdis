package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type Signable interface {
	Hash() ([32]byte, error)
	SignWith(signature []byte) error
	Signature() ([]byte, error)
}

func Sign(signable Signable, privateKey *rsa.PrivateKey) error {
	hash, err := signable.Hash()
	if err != nil {
		return err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return err
	}
	signable.SignWith([]byte(signature))
	return nil
}

func Verify(signable Signable, publicKey *rsa.PublicKey) error {
	hash, err := signable.Hash()
	if err != nil {
		return err
	}

	signature, err := signable.Signature()
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}

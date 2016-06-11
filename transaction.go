package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

type Transaction struct {
	signature []byte
	Header    TransactionHeader
}

type TransactionHeader struct {
	From  string
	To    string
	What  string
	Time  time.Time
	Nonce uint64
}

func (t *Transaction) Hash() ([32]byte, error) {
	data, err := json.Marshal(t.Header)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(data), nil
}

func (t *Transaction) NextTry() {
	t.Header.Nonce++
}

func (t *Transaction) SignWith(signature []byte) error {
	t.signature = []byte(base64.StdEncoding.EncodeToString(signature))
	return nil
}

func (t *Transaction) Signature() ([]byte, error) {
	signature, err := base64.StdEncoding.DecodeString(string(t.signature))
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (t *Transaction) Command() (Command, error) {
	var cmd Command
	err := json.Unmarshal([]byte(t.Header.What), &cmd)
	if err != nil {
		return Command{}, err
	}

	return cmd, nil
}

func NewTransaction(from string, to string, what string) *Transaction {
	return &Transaction{
		Header: TransactionHeader{
			From: from,
			To:   to,
			What: what,
			Time: time.Now(),
		},
	}
}

func NewTransactionFromCommand(from string, command Command) (*Transaction, error) {
	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	return NewTransaction(from, command.Key, string(payload)), nil
}

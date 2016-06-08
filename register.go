package main

// LWW register
import "encoding/json"

type opRegister struct {
	Key   string
	Value []byte
	Type  string
}

func NewRegister(key string, value []byte, owner string) (*Transaction, error) {
	op := opRegister{key, value, "set_register"}
	payload, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	return NewTransaction(owner, "global", string(payload)), nil
}

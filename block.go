package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Block struct {
	Header       BlockHeader
	Transactions []*Transaction
	signature    []byte
	State        State
	Previous     *Block
}

type BlockHeader struct {
	Prev     [32]byte
	RootHash [32]byte // TODO: root of merkel tree
	Time     time.Time
	Nonce    uint64
}

type State map[string]*Value

type Value struct {
	Val        interface{}
	Expire     time.Time
	WillExpire bool
}

func (v *Value) UpdateVal(val interface{}) {
	v.Val = val
}

func (v *Value) UpdateExpire(expire time.Time) {
	v.Expire = expire
	v.WillExpire = true
}

func (b *Block) Hash() ([32]byte, error) {
	data, err := json.Marshal(b.Header)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(data), nil
}

func (b *Block) NextTry() {
	b.Header.Nonce++
}

func (b *Block) SignWith(signature []byte) error {
	b.signature = []byte(base64.StdEncoding.EncodeToString(signature))
	return nil
}

func (b *Block) Signature() ([]byte, error) {
	signature, err := base64.StdEncoding.DecodeString(string(b.signature))
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

// merkle root
func (b *Block) HashTransactions() error {
	if len(b.Transactions) < 2 || !isPowerOf2(len(b.Transactions)) {
		return errors.New("block can only contain 2^n transactions and n can't be 0")
	}

	length := len(b.Transactions)
	rootHash, err := merkleHash(b.Transactions[0:length/2], b.Transactions[length/2:length])
	if err != nil {
		return err
	}
	b.Header.RootHash = rootHash
	return nil
}

func (b *Block) Verify() error {
	reached, err := reachThreshold(b)
	if err != nil {
		return err
	}
	if !reached {
		hash, err := b.Hash()
		if err != nil {
			return err
		}
		return fmt.Errorf("Invalid Proof of work on block %s", readableHash(hash))
	}

	return b.VerifyTransactions()
}

func (b *Block) VerifyTransactions() error {
	if len(b.Transactions) < 2 || !isPowerOf2(len(b.Transactions)) {
		return errors.New("block can only contain 2^n transactions and n can't be 0")
	}

	length := len(b.Transactions)
	rootHash, err := merkleHash(b.Transactions[0:length/2], b.Transactions[length/2:length])
	if err != nil {
		return err
	}
	if rootHash != b.Header.RootHash {
		return errors.New("Verification failed")
	}

	for _, tx := range b.Transactions {
		hash, err := tx.ReadableHash()
		if err != nil {

		}
		reached, err := reachThreshold(tx)
		if err != nil {
			return err
		}
		if !reached {
			return fmt.Errorf("Invalid Proof of Work %s", string(hash[:]))
		}
	}
	return nil
}

func (b *Block) UpdateState() error {
	var state State
	if b.Previous == nil {
		state = State{}
	} else {
		state = cloneState(b.Previous.State)
	}

	for _, tx := range b.Transactions {
		cmd, err := tx.Command()
		if err != nil {
			return err
		}

		// TODO: handle return values from transaction (set to another key?)
		ret, err := cmd.Execute(state)
		if err != nil {
			return err
		}

		retKey, err := tx.ReadableHash()
		if err != nil {
			return err
		}
		state[string(retKey)+":ret"] = &Value{Val: ret}
	}

	// expire expired values
	for k, v := range state {
		if v.WillExpire && time.Now().After(v.Expire) {
			delete(state, k)
		}
	}

	// TODO: hash states in blockchain with patricia tree
	b.State = state

	return nil
}

func NewBlock(previous *Block) (*Block, error) {
	var prevHash [32]byte
	var err error
	if previous != nil {
		prevHash, err = previous.Hash()
		if err != nil {
			return nil, err
		}
	}
	return &Block{
		Transactions: make([]*Transaction, 0),
		Previous:     previous,
		Header: BlockHeader{
			Time: time.Now(),
			Prev: prevHash,
		},
	}, nil
}

func isPowerOf2(n int) bool {
	return ((n & (n - 1)) == 0)
}

func merkleHash(left []*Transaction, right []*Transaction) ([32]byte, error) {
	if len(left) == 1 && len(right) == 1 {
		leftHash, err := left[0].Hash()
		if err != nil {
			return [32]byte{}, err
		}
		rightHash, err := right[0].Hash()
		if err != nil {
			return [32]byte{}, err
		}
		combine := []byte{}
		combine = append(combine, leftHash[:]...)
		combine = append(combine, rightHash[:]...)
		return sha256.Sum256(combine), nil
	}

	combine := []byte{}
	leftMerkleHash, err := merkleHash(left[0:len(left)/2], left[len(left)/2:len(left)])
	if err != nil {
		return [32]byte{}, err
	}
	rightMerkleHash, err := merkleHash(right[0:len(right)/2], right[len(right)/2:len(right)])
	if err != nil {
		return [32]byte{}, err
	}
	combine = append(combine, leftMerkleHash[:]...)
	combine = append(combine, rightMerkleHash[:]...)
	return sha256.Sum256(combine), nil
}

func cloneState(state State) State {
	newState := State{}
	for k, v := range state {
		newState[k] = v
	}

	return newState
}

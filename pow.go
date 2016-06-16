package main

import "encoding/binary"

type Workable interface {
	Hash() ([32]byte, error)
	NextTry()
}

var ProofOfWorkThreshold = [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func Work(workable Workable) error {
	for {
		reached, err := reachThreshold(workable)
		if err != nil {
			return err
		}
		if reached {
			break
		}

		workable.NextTry()
	}

	return nil
}

func reachThreshold(workable Workable) (bool, error) {
	hash, err := workable.Hash()
	if err != nil {
		return false, err
	}

	return binary.BigEndian.Uint64(hash[:]) < binary.BigEndian.Uint64(ProofOfWorkThreshold[:]), nil
}

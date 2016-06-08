package main

import "encoding/binary"

type Workable interface {
	Hash() ([32]byte, error)
	NextTry()
}

var ProofOfWorkThreshold = [32]byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func Work(workable Workable) error {
	for {
		hash, err := workable.Hash()
		if err != nil {
			return err
		}

		if binary.BigEndian.Uint64(hash[:]) < binary.BigEndian.Uint64(ProofOfWorkThreshold[:]) {
			break
		}

		workable.NextTry()
	}

	return nil
}

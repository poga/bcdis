package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	Convey("INCR", t, func() {
		previousBlock, err := setupPreviousBlockState(map[string]interface{}{"foo": "1"})
		So(err, ShouldBeNil)

		Convey("can update state", func() {
			block, err := setupTestCommandBlock(previousBlock, NewCommand(INCR, "foo"))
			So(err, ShouldBeNil)

			So(block.State, ShouldResemble, map[string]interface{}{"foo": "2"})
		})

		Convey("initialize key as 0 if not exists", func() {
			block, err := setupTestCommandBlock(previousBlock, NewCommand(INCR, "bar"))
			So(err, ShouldBeNil)

			So(block.State, ShouldResemble, map[string]interface{}{"foo": "1", "bar": "1"})
		})

		Convey("return err if value can't be parsed as string", func() {
			previousBlock, err := setupPreviousBlockState(map[string]interface{}{"foo": "x"})
			So(err, ShouldBeNil)

			_, err = setupTestCommandBlock(previousBlock, NewCommand(INCR, "foo"))
			So(err, ShouldNotBeNil)
		})

	})
}

func setupPreviousBlockState(previousState map[string]interface{}) (*Block, error) {
	block, err := NewBlock(nil)
	if err != nil {
		return nil, err
	}

	// setup previous block state
	for k, v := range previousState {
		tx, err := NewTransactionFromCommand("alice", NewCommand(SET, k, v.(string)))
		if err != nil {
			return nil, err
		}

		block.Transactions = append(block.Transactions, tx)
	}

	err = block.UpdateState()
	if err != nil {
		return nil, err
	}
	return block, nil
}

func setupTestCommandBlock(previousBlock *Block, cmd Command) (*Block, error) {
	block, err := NewBlock(previousBlock)
	if err != nil {
		return nil, err
	}

	// setup previous block state
	tx, err := NewTransactionFromCommand("alice", cmd)
	if err != nil {
		return nil, err
	}

	block.Transactions = append(block.Transactions, tx)

	err = block.UpdateState()
	if err != nil {
		return nil, err
	}
	return block, nil
}

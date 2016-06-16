package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBlock(t *testing.T) {
	Convey("A block", t, func() {
		rootBlock, err := NewBlock(nil)
		So(err, ShouldBeNil)
		childBlock, err := NewBlock(rootBlock)
		So(err, ShouldBeNil)

		testblocks := map[string]*Block{
			"root block":  rootBlock,
			"child block": childBlock,
		}

		Convey("child block knows the hash of parent block", func() {
			rootHash, err := rootBlock.Hash()
			So(err, ShouldBeNil)
			So(rootBlock.Header.Prev, ShouldEqual, [32]byte{})
			So(childBlock.Header.Prev, ShouldEqual, rootHash)
		})

		for name, b := range testblocks {

			Convey(name+" is hashable", func() {
				_, err := b.Hash()
				So(err, ShouldBeNil)
			})

			Convey(name+" is workable", func() {
				Convey("can have muliple try of Proof of Work", func() {
					hash, err := b.Hash()
					So(err, ShouldBeNil)

					b.NextTry()

					hash2, err := b.Hash()
					So(err, ShouldBeNil)
					So(hash, ShouldNotEqual, hash2)
				})

				Convey(" can find a correct proof of work", func() { // might take some time
					err := Work(b)

					So(err, ShouldBeNil)
				})
			})

			Convey(name+" is signable", func() {
				account, err := NewAccount()
				So(err, ShouldBeNil)

				Convey("can be signed without change it's hash", func() {
					hash, err := b.Hash()
					So(err, ShouldBeNil)

					Sign(b, account)

					hash2, err := b.Hash()
					So(err, ShouldBeNil)
					So(hash, ShouldEqual, hash2)
				})

				Convey("can be verified via its signature", func() {
					hash, err := b.Hash()
					So(err, ShouldBeNil)

					Sign(b, account)

					hash2, err := b.Hash()
					So(err, ShouldBeNil)
					So(hash, ShouldEqual, hash2)

					err = Verify(b, account)
					So(err, ShouldBeNil)
				})
				Convey("can be signed without affecting hash", func() {
					hash, err := b.Hash()
					So(err, ShouldBeNil)

					b.SignWith([]byte("foo"))

					hash2, err := b.Hash()
					So(err, ShouldBeNil)
					So(hash, ShouldEqual, hash2)
				})
			})

			Convey(name+" with no transaction can't be verified", func() {
				So(b.HashTransactions(), ShouldNotBeNil)
				So(b.VerifyTransactions(), ShouldNotBeNil)
			})

			Convey(name+" can add transactions into block", func() {
				tx := NewTransaction("alice", "bob", "payload")
				// find valid proof of work
				So(Work(tx), ShouldBeNil)
				b.Transactions = append(b.Transactions, tx)

				Convey("block transaction count need to be power of 2 to calculate merkle hash", func() {
					So(b.HashTransactions(), ShouldNotBeNil)
					So(b.VerifyTransactions(), ShouldNotBeNil)
				})

				Convey("block with 2 transaction can be hashed with merkle tree", func() {
					tx := NewTransaction("bob", "alice", "payload2")
					So(Work(tx), ShouldBeNil)

					b.Transactions = append(b.Transactions, tx)
					So(b.HashTransactions(), ShouldBeNil)
					So(b.Header.RootHash, ShouldNotEqual, [32]byte{})
					So(b.VerifyTransactions(), ShouldBeNil)

					Convey("if any transaction is modified, block should be failed to verify without rehash", func() {
						b.Transactions[0].NextTry()
						So(b.VerifyTransactions(), ShouldNotBeNil)

						Convey("if we rehash the block and rework the transaction in the block, the block will be verified again", func() {
							So(Work(b.Transactions[0]), ShouldBeNil)
							So(b.HashTransactions(), ShouldBeNil)
							So(b.Header.RootHash, ShouldNotEqual, [32]byte{})
							So(b.VerifyTransactions(), ShouldBeNil)

							Convey("block can be worked to be valid", func() {
								So(Work(b), ShouldBeNil)
								So(b.Verify(), ShouldBeNil)
							})
						})
					})

					Convey("block can be worked to be valid", func() {
						So(Work(b), ShouldBeNil)
						So(b.Verify(), ShouldBeNil)
					})
				})

				Convey("block with len(transaction) = 2^n can be hashed with merkle tree", func() {
					tx := NewTransaction("bob", "alice", "payload2")
					So(Work(tx), ShouldBeNil)
					b.Transactions = append(b.Transactions, tx)

					tx = NewTransaction("alice", "bob", "payload3")
					So(Work(tx), ShouldBeNil)
					b.Transactions = append(b.Transactions, tx)

					tx = NewTransaction("bob", "alice", "payload4")
					So(Work(tx), ShouldBeNil)
					b.Transactions = append(b.Transactions, tx)

					So(b.HashTransactions(), ShouldBeNil)
					So(b.Header.RootHash, ShouldNotEqual, [32]byte{})
					So(b.VerifyTransactions(), ShouldBeNil)

					Convey("block can be worked to be valid", func() {
						So(Work(b), ShouldBeNil)
						So(b.Verify(), ShouldBeNil)
					})

					Convey("if any transaction is modified, block should be failed to verify without rehash", func() {
						b.Transactions[1].NextTry()
						So(b.VerifyTransactions(), ShouldNotBeNil)

						Convey("if we rehash the block and rework the transaction, the block will be verified again", func() {
							So(Work(b.Transactions[1]), ShouldBeNil)

							So(b.HashTransactions(), ShouldBeNil)
							So(b.Header.RootHash, ShouldNotEqual, [32]byte{})
							So(b.VerifyTransactions(), ShouldBeNil)

							Convey("block can be worked to be valid", func() {
								So(Work(b), ShouldBeNil)
								So(b.Verify(), ShouldBeNil)
							})
						})
					})
				})
			})
		}
	})

	Convey("A root block with commmand transaction", t, func() {
		rootBlock, err := NewBlock(nil)
		So(err, ShouldBeNil)

		tx, err := NewTransactionFromCommand("alice", NewCommand(SET, "foo", "bar"))
		So(err, ShouldBeNil)

		rootBlock.Transactions = append(rootBlock.Transactions, tx)

		Convey("can caluclate new state based on included transactions", func() {
			err := rootBlock.UpdateState()
			So(err, ShouldBeNil)

			So(rootBlock.State["foo"].Val, ShouldEqual, "bar")

			Convey("A child block with command transaction", func() {
				childBlock, err := NewBlock(rootBlock)
				So(err, ShouldBeNil)

				tx, err := NewTransactionFromCommand("alice", NewCommand(SET, "foo2", "baz"))
				So(err, ShouldBeNil)

				Convey("can caluclate new state based on included transactions", func() {
					childBlock.Transactions = append(childBlock.Transactions, tx)
					err := childBlock.UpdateState()
					So(err, ShouldBeNil)

					// previous state should not be affected
					So(rootBlock.State["foo"].Val, ShouldEqual, "bar")

					So(childBlock.State["foo2"].Val, ShouldEqual, "baz")
				})
			})

			Convey("can set return value from transaction into another key", func() {
				childBlock, err := NewBlock(rootBlock)
				So(err, ShouldBeNil)

				tx, err := NewTransactionFromCommand("alice", NewCommand(GETSET, "foo", "baz"))
				So(err, ShouldBeNil)

				Convey("can caluclate new state based on included transactions", func() {
					childBlock.Transactions = append(childBlock.Transactions, tx)
					err := childBlock.UpdateState()
					So(err, ShouldBeNil)

					// previous state should not be affected
					So(rootBlock.State["foo"].Val, ShouldEqual, "bar")
					retKey, err := tx.ReadableHash()
					So(err, ShouldBeNil)

					So(childBlock.State["foo"].Val, ShouldEqual, "baz")
					So(childBlock.State[string(retKey)+":ret"].Val, ShouldEqual, "bar")
				})
			})
		})

		Convey("can expire value when updating state", func() {
			err := rootBlock.UpdateState()
			So(err, ShouldBeNil)

			So(rootBlock.State["foo"].Val, ShouldEqual, "bar")

			Convey("A child block with command transaction", func() {
				childBlock, err := NewBlock(rootBlock)
				So(err, ShouldBeNil)

				tx, err := NewTransactionFromCommand("alice", NewCommand(SET, "foo2", "baz"))
				So(err, ShouldBeNil)

				Convey("can caluclate new state based on included transactions", func() {
					childBlock.Transactions = append(childBlock.Transactions, tx)
					err := childBlock.UpdateState()
					So(err, ShouldBeNil)

					// previous state should not be affected
					So(rootBlock.State["foo"].Val, ShouldEqual, "bar")

					So(childBlock.State["foo2"].Val, ShouldEqual, "baz")
				})
			})

			Convey("can set return value from transaction into another key", func() {
				childBlock, err := NewBlock(rootBlock)
				So(err, ShouldBeNil)

				tx, err := NewTransactionFromCommand("alice", NewCommand(GETSET, "foo", "baz"))
				So(err, ShouldBeNil)

				Convey("can caluclate new state based on included transactions", func() {
					childBlock.Transactions = append(childBlock.Transactions, tx)
					err := childBlock.UpdateState()
					So(err, ShouldBeNil)

					// previous state should not be affected
					So(rootBlock.State["foo"].Val, ShouldEqual, "bar")
					retKey, err := tx.ReadableHash()
					So(err, ShouldBeNil)

					So(childBlock.State["foo"].Val, ShouldEqual, "baz")
					So(childBlock.State[string(retKey)+":ret"].Val, ShouldEqual, "bar")
				})
			})
		})

		Convey("A root block with expire transaction", func() {
			rootBlock, err := NewBlock(nil)
			So(err, ShouldBeNil)

			tx, err := NewTransactionFromCommand("alice", NewCommand(SET, "foo", "bar"))
			So(err, ShouldBeNil)
			rootBlock.Transactions = append(rootBlock.Transactions, tx)

			tx, err = NewTransactionFromCommand("alice", NewCommand(EXPIRE, "foo", "2"))
			So(err, ShouldBeNil)
			rootBlock.Transactions = append(rootBlock.Transactions, tx)

			time.Sleep(3 * time.Second)

			Convey("can expire value when updating state", func() {
				err := rootBlock.UpdateState()
				So(err, ShouldBeNil)

				So(rootBlock.State["foo"], ShouldBeNil)

			})
		})
	})
}

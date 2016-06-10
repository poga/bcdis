package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	Convey("A transaction", t, func() {
		t := NewTransaction("alice", "bob", "op")

		Convey("is hashable", func() {
			_, err := t.Hash()
			So(err, ShouldBeNil)
		})

		Convey("is workable", func() {
			Convey("can have muliple try of Proof of Work", func() {
				hash, err := t.Hash()
				So(err, ShouldBeNil)

				t.NextTry()

				hash2, err := t.Hash()
				So(err, ShouldBeNil)
				So(hash, ShouldNotEqual, hash2)
			})

			Convey("can find a correct proof of work", func() { // might take some time
				err := Work(t)

				So(err, ShouldBeNil)
			})
		})

		Convey("is signable", func() {
			account, err := NewAccount()
			So(err, ShouldBeNil)

			Convey("can be signed without change it's hash", func() {
				hash, err := t.Hash()
				So(err, ShouldBeNil)

				Sign(t, account)

				hash2, err := t.Hash()
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, hash2)
			})

			Convey("can be verified via its signature", func() {
				hash, err := t.Hash()
				So(err, ShouldBeNil)

				Sign(t, account)

				hash2, err := t.Hash()
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, hash2)

				err = Verify(t, account)
				So(err, ShouldBeNil)
			})
			Convey("can be signed without affecting hash", func() {
				hash, err := t.Hash()
				So(err, ShouldBeNil)

				t.SignWith([]byte("foo"))

				hash2, err := t.Hash()
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, hash2)
			})
		})

	})
}

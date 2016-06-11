package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAccount(t *testing.T) {
	Convey("an account", t, func() {
		a, err := NewAccount()
		So(err, ShouldBeNil)

		Convey("have private part and public part", func() {
			// public key
			So(a.Public(), ShouldNotBeNil)

			// private key
			So(a.Key, ShouldNotBeNil)
		})

		Convey("have an human readable address", func() {
			address, err := a.Address()
			So(err, ShouldBeNil)

			So(string(address[:]), ShouldNotBeNil)
		})
	})
}

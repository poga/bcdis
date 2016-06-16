package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	Convey("INCR", t, func() {
		state := State{"foo": &Value{Val: "1"}}

		Convey("can update state", func() {
			command := NewCommand(INCR, "foo")

			_, err := command.Execute(state)
			So(err, ShouldBeNil)
			So(state["foo"].Val, ShouldEqual, "2")
		})

		Convey("initialize key as 0 if not exists", func() {
			command := NewCommand(INCR, "bar")
			_, err := command.Execute(state)

			So(err, ShouldBeNil)
			So(state["foo"].Val, ShouldEqual, "1")
			So(state["bar"].Val, ShouldEqual, "1")
		})

		Convey("return err if value can't be parsed as string", func() {
			state = State{"foo": &Value{Val: "x"}}
			command := NewCommand(INCR, "foo")

			_, err := command.Execute(state)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("GET", t, func() {
		state := State{"foo": &Value{Val: "1"}}

		Convey("can return state", func() {
			cmd := NewCommand(GET, "foo")

			ret, err := cmd.Execute(state)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "1")

			So(state["foo"].Val, ShouldEqual, "1")
		})
	})

	Convey("GETSET", t, func() {
		state := State{"foo": &Value{Val: "1"}}

		Convey("can update and set state", func() {
			cmd := NewCommand(GETSET, "foo", "2")

			ret, err := cmd.Execute(state)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "1")

			So(state["foo"].Val, ShouldEqual, "2")
		})

		Convey("returns error if key is associated with a non-string value", func() {
			state := State{"foo": &Value{Val: 1}}

			cmd := NewCommand(GETSET, "foo", "2")

			_, err := cmd.Execute(state)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("EXPIRE", t, func() {
		state := State{"foo": &Value{Val: "1"}}

		Convey("can set a value to correct expire time", func() {
			cmd := NewCommand(EXPIRE, "foo", "1")
			// expire command need to be in a transaction to know when to expire
			cmd.TX = &Transaction{Header: TransactionHeader{Time: time.Now()}}

			ret, err := cmd.Execute(state)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "OK")

			So(state["foo"].Val, ShouldEqual, "1")
			So(state["foo"].WillExpire, ShouldEqual, true)
		})
	})
}

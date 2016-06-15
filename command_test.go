package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	Convey("INCR", t, func() {
		state := map[string]interface{}{"foo": "1"}

		Convey("can update state", func() {
			command := NewCommand(INCR, "foo")

			_, err := command.Execute(state)
			So(err, ShouldBeNil)
			So(state, ShouldResemble, map[string]interface{}{"foo": "2"})
		})

		Convey("initialize key as 0 if not exists", func() {
			command := NewCommand(INCR, "bar")
			_, err := command.Execute(state)

			So(err, ShouldBeNil)
			So(state, ShouldResemble, map[string]interface{}{"foo": "1", "bar": "1"})
		})

		Convey("return err if value can't be parsed as string", func() {
			state = map[string]interface{}{"foo": "x"}
			command := NewCommand(INCR, "foo")

			_, err := command.Execute(state)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("GET", t, func() {
		state := map[string]interface{}{"foo": "1"}

		Convey("can return state", func() {
			cmd := NewCommand(GET, "foo")

			ret, err := cmd.Execute(state)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, "1")

			So(state, ShouldResemble, map[string]interface{}{"foo": "1"})
		})
	})
}

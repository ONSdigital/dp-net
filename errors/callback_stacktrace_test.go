package errors_test

import (
	"fmt"
	"testing"

	dperrors "github.com/ONSdigital/dp-net/v2/errors"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	packagePath = "dp-net/errors"
	fileName    = "callback_stacktrace_test.go"
)

// These tests should stay in their own file as they rely on specific line numbers of
// their own code, having them in a file that will be continuously changing would
// introduce unnecessary overhead to maintain
func TestStackTraceHappy(t *testing.T) {
	Convey("Given an error with embedded stack trace from pkg/errors", t, func() {
		err := testCallStackFunc1()
		Convey("When stackTrace(err) is called", func() {
			st := dperrors.StackTrace(err)
			So(len(st), ShouldEqual, 19)

			So(st[0].File, ShouldContainSubstring, packagePath+"/"+fileName)
			So(st[0].Line, ShouldEqual, 70)
			So(st[0].Function, ShouldEqual, "testCallStackFunc3")

			So(st[1].File, ShouldContainSubstring, packagePath+"/"+fileName)
			So(st[1].Line, ShouldEqual, 66)
			So(st[1].Function, ShouldEqual, "testCallStackFunc2")

			So(st[2].File, ShouldContainSubstring, packagePath+"/"+fileName)
			So(st[2].Line, ShouldEqual, 62)
			So(st[2].Function, ShouldEqual, "testCallStackFunc1")
		})
	})

	Convey("Given an error with intermittently embedded stack traces from pkg/errors", t, func() {
		err := testCallStackFunc4()
		Convey("When stackTrace(err) is called", func() {
			st := dperrors.StackTrace(err)
			So(len(st), ShouldEqual, 18)

			So(st[0].File, ShouldContainSubstring, packagePath+"/"+fileName)
			So(st[0].Line, ShouldEqual, 70)
			So(st[0].Function, ShouldEqual, "testCallStackFunc3")

			So(st[1].File, ShouldContainSubstring, packagePath+"/"+fileName)
			So(st[1].Line, ShouldEqual, 75)
			So(st[1].Function, ShouldEqual, "testCallStackFunc4")

		})
	})
}

func testCallStackFunc1() error {
	return testCallStackFunc2()
}

func testCallStackFunc2() error {
	return testCallStackFunc3()
}

func testCallStackFunc3() error {
	cause := errors.New("I am the cause")
	return errors.Wrap(cause, "I am the context")
}

func testCallStackFunc4() error {
	if err := testCallStackFunc3(); err != nil {
		return fmt.Errorf("I do not have embedded stack trace, but this cause does: %w", err)
	}
	return nil
}

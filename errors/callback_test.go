package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type testError struct {
	err        error
	statusCode int
	logData    map[string]interface{}
}

func (e testError) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

func (e testError) Unwrap() error {
	return e.err
}

func (e testError) Code() int {
	return e.statusCode
}

func (e testError) LogData() map[string]interface{} {
	return e.logData
}

func TestUnwrapLogDataHappy(t *testing.T) {

	Convey("Given an error with embedded logData", t, func() {
		err := &testError{
			logData: log.Data{
				"log": "data",
			},
		}

		Convey("When logData(err) is called", func() {
			ld := dperrors.LogData(err)
			So(ld, ShouldResemble, log.Data{"log": "data"})
		})
	})

	Convey("Given an error chain with wrapped logData", t, func() {
		err1 := &testError{
			err: errors.New("original error"),
			logData: log.Data{
				"log": "data",
			},
		}

		err2 := &testError{
			err: fmt.Errorf("err1: %w", err1),
			logData: log.Data{
				"additional": "data",
			},
		}

		err3 := &testError{
			err: fmt.Errorf("err2: %w", err2),
			logData: log.Data{
				"final": "data",
			},
		}

		Convey("When unwrapLogData(err) is called", func() {
			logData := dperrors.UnwrapLogData(err3)
			expected := log.Data{
				"final":      "data",
				"additional": "data",
				"log":        "data",
			}

			So(logData, ShouldResemble, expected)
		})
	})

	Convey("Given an error chain with intermittent wrapped logData", t, func() {
		err1 := &testError{
			err: errors.New("original error"),
			logData: log.Data{
				"log": "data",
			},
		}

		err2 := &testError{
			err: fmt.Errorf("err1: %w", err1),
		}

		err3 := &testError{
			err: fmt.Errorf("err2: %w", err2),
			logData: log.Data{
				"final": "data",
			},
		}

		Convey("When unwrapLogData(err) is called", func() {
			logData := dperrors.UnwrapLogData(err3)
			expected := log.Data{
				"final": "data",
				"log":   "data",
			}

			So(logData, ShouldResemble, expected)
		})
	})

	Convey("Given an error chain with wrapped logData with duplicate key values", t, func() {
		err1 := &testError{
			err: errors.New("original error"),
			logData: log.Data{
				"log":        "data",
				"duplicate":  "duplicate_data1",
				"request_id": "ADB45F",
			},
		}

		err2 := &testError{
			err: fmt.Errorf("err1: %w", err1),
			logData: log.Data{
				"additional": "data",
				"duplicate":  "duplicate_data2",
				"request_id": "ADB45F",
			},
		}

		err3 := &testError{
			err: fmt.Errorf("err2: %w", err2),
			logData: log.Data{
				"final":      "data",
				"duplicate":  "duplicate_data3",
				"request_id": "ADB45F",
			},
		}

		Convey("When unwrapLogData(err) is called", func() {
			logData := dperrors.UnwrapLogData(err3)
			expected := log.Data{
				"final":      "data",
				"additional": "data",
				"log":        "data",
				"duplicate": []interface{}{
					"duplicate_data3",
					"duplicate_data2",
					"duplicate_data1",
				},
				"request_id": "ADB45F",
			}

			So(logData, ShouldResemble, expected)
		})
	})
}

func TestUnwrapStatusCodeHappy(t *testing.T) {

	Convey("Given an error with embedded status code", t, func() {
		err := &testError{
			statusCode: http.StatusTeapot,
		}

		Convey("When StatusCode(err) is called", func() {
			status := dperrors.StatusCode(err)
			expected := http.StatusTeapot

			So(status, ShouldEqual, expected)
		})
	})

	Convey("Given an error chain with embedded status code", t, func() {
		err1 := &testError{
			err:        errors.New("original error"),
			statusCode: http.StatusTooManyRequests,
		}

		err2 := &testError{
			err: fmt.Errorf("err1: %w", err1),
		}

		err3 := &testError{
			err: fmt.Errorf("err2: %w", err2),
		}

		Convey("When UnwrapStatusCode(err) is called", func() {
			status := dperrors.UnwrapStatusCode(err3)
			expected := http.StatusTooManyRequests

			So(status, ShouldEqual, expected)
		})
	})

	Convey("Given an error chain with multiple embedded status codes", t, func() {
		err1 := &testError{
			err:        errors.New("original error"),
			statusCode: http.StatusBadRequest,
		}

		err2 := &testError{
			err:        fmt.Errorf("err1: %w", err1),
			statusCode: http.StatusUnauthorized,
		}

		err3 := &testError{
			err: fmt.Errorf("err2: %w", err2),
		}

		Convey("When UnwrapStatusCode(err) is called", func() {
			status := dperrors.UnwrapStatusCode(err3)
			expected := http.StatusUnauthorized
			Convey("The first valid status code is returned ", func() {

				So(status, ShouldEqual, expected)
			})
		})
	})

	Convey("Given an error with no embedded status code", t, func() {
		err := &testError{}

		Convey("When StatusCode(err) is called", func() {
			status := dperrors.UnwrapStatusCode(err)
			expected := http.StatusInternalServerError

			So(status, ShouldEqual, expected)
		})
	})

}

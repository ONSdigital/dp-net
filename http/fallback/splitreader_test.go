package fallback

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadCloserSplit_New(t *testing.T) {
	Convey("Given there is an upstream ReadCloser", t, func() {
		readCloser := io.NopCloser(bytes.NewReader([]byte("hello world")))
		Convey("When it is split into multiple readers", func() {
			splitRCs := ReadCloserSplit(readCloser, 3)
			Convey("Then a slice of readers is returned with the correct size ", func() {
				So(splitRCs, ShouldHaveLength, 3)
			})
		})
	})
}

func TestReadCloserSplit_Read(t *testing.T) {
	Convey("Given there is an upstream ReadCloser with unread bytes split into 1 reader", t, func() {
		readCloser := io.NopCloser(bytes.NewReader([]byte("hello world")))
		splitRCs := ReadCloserSplit(readCloser, 1)
		rc1 := splitRCs[0]
		Convey("When the reader reads a number of bytes", func() {
			data := make([]byte, 2)
			n, err := rc1.Read(data)
			Convey("Then it returns without error", func() {
				So(n, ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given there is an upstream ReadCloser with no unread bytes split into 1 reader", t, func() {
		readCloser := io.NopCloser(bytes.NewReader([]byte{}))
		splitRCs := ReadCloserSplit(readCloser, 1)
		rc1 := splitRCs[0]
		Convey("When the reader reads a number of bytes", func() {
			data := make([]byte, 2)
			n, err := rc1.Read(data)
			Convey("Then it returns an EOF error", func() {
				So(n, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "EOF")
			})
		})
	})
}

func TestReadCloserSplit_Close(t *testing.T) {
	Convey("Given there is an upstream ReadCloser with unread bytes split into 1 reader", t, func() {
		upstreamRC := mockReadCloser()
		splitRCs := ReadCloserSplit(upstreamRC, 1)
		rc1 := splitRCs[0]
		Convey("When the reader is closed", func() {
			err := rc1.Close()
			Convey("Then it closes without error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the upstream ReadCloser is closed without error", func() {
				So(upstreamRC.CloseCount, ShouldEqual, 1)
			})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes split into 2 readers", t, func() {
		upstreamRC := mockReadCloser()
		splitRCs := ReadCloserSplit(upstreamRC, 2)
		rc1 := splitRCs[0]
		rc2 := splitRCs[1]

		Convey("When the first reader is closed", func() {
			err := rc1.Close()
			Convey("Then it closes without error", func() {
				So(err, ShouldBeNil)
				Convey("And the upstream ReadCloser is not closed", func() {
					So(upstreamRC.CloseCount, ShouldEqual, 0)
					Convey("When the second reader is closed", func() {
						err := rc2.Close()
						Convey("Then it closes without errors", func() {
							So(err, ShouldBeNil)
							Convey("And the upstream ReadCloser is closed without errors", func() {
								So(upstreamRC.CloseCount, ShouldEqual, 1)
							})
						})
					})
				})
			})
		})
	})

}

func TestReadCloserSplit(t *testing.T) {
	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		Convey("And it is split into two readers", func() {
			Convey("When the first reader reads a number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the second reader reads the same number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("And the upstream ReadCloser should have been read from only once", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		Convey("And it is split into two readers", func() {
			Convey("When the first reader reads a number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the second reader reads a greater number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the first reader reads up to the same index as the second reader", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("And the upstream ReadCloser should have been read from twice", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		Convey("And it is split into two readers", func() {
			Convey("When the first reader reads a number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the second reader reads a greater number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the first reader reads beyond the index of the second reader", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("And the upstream ReadCloser should have been read from three times", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		Convey("And it is split into two readers", func() {
			Convey("When the first reader reads a number of bytes", func() {})
			Convey("And the first reader is then closed", func() {})
			Convey("When the second reader reads a greater number of bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("And the upstream ReadCloser should have been read from twice", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser which errors after 5 bytes are read", t, func() {
		Convey("And it is split into three readers", func() {
			Convey("When the first reader reads 5 bytes", func() {
				Convey("Then it retuns the bytes correctly", func() {})
			})
			Convey("When the first reader reads one more byte", func() {
				Convey("Then it retuns an error", func() {})
			})
			Convey("When the second reader reads 5 bytes", func() {
				Convey("Then it return the bytes correctly", func() {})
			})
			Convey("When the second reader reads one more byte", func() {
				Convey("Then it retuns an error", func() {})
			})
			Convey("When the third reader reads 6 bytes", func() {
				Convey("Then it retuns an error", func() {})
			})
		})
	})
}

// Mock for io.ReadCloser
type readCloserMock struct {
	CloseCount int
}

func (rcm *readCloserMock) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (rcm *readCloserMock) Close() error {
	rcm.CloseCount++
	return nil
}

var _ io.ReadCloser = &readCloserMock{}

func mockReadCloser() *readCloserMock {
	return &readCloserMock{}
}

func (rcm *readCloserMock) reset() {
	rcm.CloseCount = 0
}

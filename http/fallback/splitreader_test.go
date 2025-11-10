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
		upstreamRC := mockReadCloser(nil)
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
		upstreamRC := mockReadCloser(nil)
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
		upstreamRC := mockReadCloser([]byte("some.data.to.test"))
		Convey("And it is split into two readers", func() {
			splitRCs := ReadCloserSplit(upstreamRC, 2)
			rc1 := splitRCs[0]
			rc2 := splitRCs[1]
			Convey("When the first reader reads a number of bytes", func() {
				read1 := make([]byte, 5)
				n, err := rc1.Read(read1)
				So(err, ShouldBeNil)
				Convey("Then it return the bytes correctly", func() {
					So(n, ShouldEqual, 5)
					So(read1, ShouldResemble, []byte("some."))
					Convey("When the second reader reads the same number of bytes", func() {
						read2 := make([]byte, 5)
						n, err = rc2.Read(read2)
						So(err, ShouldBeNil)
						Convey("Then it return the bytes correctly", func() {
							So(n, ShouldEqual, 5)
							So(read1, ShouldResemble, []byte("some."))
							Convey("And the upstream ReadCloser should have been read from only once", func() {
								So(upstreamRC.Reads, ShouldHaveLength, 1)
							})
						})
					})
				})
			})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		upstreamRC := mockReadCloser([]byte("some.data.to.test"))
		Convey("And it is split into two readers", func() {
			splitRCs := ReadCloserSplit(upstreamRC, 2)
			rc1 := splitRCs[0]
			rc2 := splitRCs[1]
			Convey("When the first reader reads a number of bytes", func() {
				read1 := make([]byte, 5)
				n, err := rc1.Read(read1)
				Convey("Then it return the bytes correctly", func() {
					So(err, ShouldBeNil)
					So(n, ShouldEqual, 5)
					So(read1, ShouldResemble, []byte("some."))
					Convey("When the second reader reads a greater number of bytes", func() {
						read2 := make([]byte, 8)
						n, err = rc2.Read(read2)
						Convey("Then it return the bytes correctly", func() {
							So(err, ShouldBeNil)
							So(n, ShouldEqual, 8)
							So(read2, ShouldResemble, []byte("some.dat"))
							Convey("When the first reader reads up to the same index as the second reader", func() {
								read1b := make([]byte, 3)
								n, err = rc1.Read(read1b)
								Convey("Then it return the bytes correctly", func() {
									So(err, ShouldBeNil)
									So(n, ShouldEqual, 3)
									So(read1b, ShouldResemble, []byte("dat"))
									Convey("And the upstream ReadCloser should have been read from twice", func() {
										So(upstreamRC.Reads, ShouldHaveLength, 2)
									})
								})
							})
						})
					})
				})
			})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		upstreamRC := mockReadCloser([]byte("some.data.to.test"))
		Convey("And it is split into two readers", func() {
			splitRCs := ReadCloserSplit(upstreamRC, 2)
			rc1 := splitRCs[0]
			rc2 := splitRCs[1]
			Convey("When the first reader reads a number of bytes", func() {
				read1 := make([]byte, 5)
				n, err := rc1.Read(read1)
				Convey("Then it return the bytes correctly", func() {
					So(err, ShouldBeNil)
					So(n, ShouldEqual, 5)
					So(read1, ShouldResemble, []byte("some."))

					Convey("When the second reader reads a greater number of bytes", func() {
						read2 := make([]byte, 8)
						n, err = rc2.Read(read2)
						Convey("Then it return the bytes correctly", func() {
							So(err, ShouldBeNil)
							So(n, ShouldEqual, 8)
							So(read2, ShouldResemble, []byte("some.dat"))

							Convey("When the first reader reads beyond the index of the second reader", func() {
								read1b := make([]byte, 6)
								n, err = rc1.Read(read1b)
								Convey("Then it return the bytes correctly", func() {
									So(err, ShouldBeNil)
									So(n, ShouldEqual, 6)
									So(read1b, ShouldResemble, []byte("data.t"))

									Convey("And the upstream ReadCloser should have been read from three times", func() {
										So(upstreamRC.Reads, ShouldHaveLength, 3)
									})
								})
							})
						})
					})
				})
			})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes", t, func() {
		upstreamRC := mockReadCloser([]byte("some.data.to.test"))
		Convey("And it is split into two readers", func() {
			splitRCs := ReadCloserSplit(upstreamRC, 2)
			rc1 := splitRCs[0]
			rc2 := splitRCs[1]
			Convey("When the first reader reads a number of bytes", func() {
				read1 := make([]byte, 5)
				n, err := rc1.Read(read1)
				So(err, ShouldBeNil)

				Convey("And the first reader is then closed", func() {
					err = rc1.Close()
					So(err, ShouldBeNil)
					Convey("When the second reader reads a greater number of bytes", func() {
						read2 := make([]byte, 8)
						n, err = rc2.Read(read2)
						Convey("Then it return the bytes correctly", func() {
							So(err, ShouldBeNil)
							So(n, ShouldEqual, 8)
							So(read2, ShouldResemble, []byte("some.dat"))

							Convey("And the upstream ReadCloser should have been read from twice", func() {
								So(upstreamRC.Reads, ShouldHaveLength, 2)
							})
						})
					})
				})
			})
		})
	})
}

func TestReadCloserSplit_Errors(t *testing.T) {
	Convey("Given there is an upstream ReadCloser which errors after 5 bytes are read", t, func() {
		upstreamRC := mockReadCloser([]byte("some."))
		Convey("And it is split into three readers", func() {
			splitRCs := ReadCloserSplit(upstreamRC, 3)
			rc1 := splitRCs[0]
			rc2 := splitRCs[1]
			rc3 := splitRCs[2]
			Convey("When the first reader reads 5 bytes", func() {
				read1 := make([]byte, 5)
				n, err := rc1.Read(read1)
				Convey("Then it returns the bytes correctly", func() {
					So(err, ShouldBeNil)
					So(n, ShouldEqual, 5)
					So(read1, ShouldResemble, []byte("some."))
					Convey("When the first reader reads one more byte", func() {
						read1b := make([]byte, 1)
						n, err := rc1.Read(read1b)
						Convey("Then it returns an error", func() {
							So(err, ShouldNotBeNil)
							So(n, ShouldEqual, 0)
							So(err, ShouldEqual, io.EOF)
							Convey("When the second reader reads 5 bytes", func() {
								read2 := make([]byte, 5)
								n, err = rc2.Read(read2)
								Convey("Then it return the bytes correctly", func() {
									So(err, ShouldBeNil)
									So(n, ShouldEqual, 5)
									So(read2, ShouldResemble, []byte("some."))
									Convey("When the second reader reads one more byte", func() {
										read2b := make([]byte, 1)
										n, err = rc2.Read(read2b)
										Convey("Then it returns an error", func() {
											So(err, ShouldNotBeNil)
											So(n, ShouldEqual, 0)
											So(err, ShouldEqual, io.EOF)
											Convey("When the third reader reads 6 bytes", func() {
												read3 := make([]byte, 6)
												n, err = rc3.Read(read3)
												Convey("Then it returns only the 5 unread bytes with no error", func() {
													So(err, ShouldBeNil)
													So(n, ShouldEqual, 5)
													Convey("And the first 5 bytes of the buffer should contain the correct data", func() {})
													So(read3[:5], ShouldResemble, []byte("some."))
													Convey("When the second reader reads one more byte", func() {
														read3b := make([]byte, 1)
														n, err = rc3.Read(read3b)
														Convey("Then it returns an error", func() {
															So(err, ShouldNotBeNil)
															So(n, ShouldEqual, 0)
															So(err, ShouldEqual, io.EOF)
														})
													})
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})
}

// Mock for io.ReadCloser
type readCloserMock struct {
	data       []byte
	reader     io.Reader
	CloseCount int
	Reads      []int
}

func (rcm *readCloserMock) Read(p []byte) (n int, err error) {
	rcm.Reads = append(rcm.Reads, len(p))
	return rcm.reader.Read(p)
}

func (rcm *readCloserMock) Close() error {
	rcm.CloseCount++
	return nil
}

var _ io.ReadCloser = &readCloserMock{}

func mockReadCloser(data []byte) *readCloserMock {
	return &readCloserMock{
		data:   data,
		reader: bytes.NewReader(data),
		Reads:  make([]int, 0),
	}
}

func (rcm *readCloserMock) reset() {
	rcm.reader = bytes.NewReader(rcm.data)
	rcm.CloseCount = 0
	rcm.Reads = make([]int, 0)
}

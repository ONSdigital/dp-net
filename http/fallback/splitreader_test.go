package fallback

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadCloserSplit_New(t *testing.T) {
	Convey("Given there is an upstream ReadCloser", t, func() {
		Convey("When it is split into multiple readers", func() {
			Convey("Then a slice of readers is returned with the correct size ", func() {})
		})
	})
}

func TestReadCloserSplit_Read(t *testing.T) {
	Convey("Given there is an upstream ReadCloser with unread bytes split into 1 reader", t, func() {
		Convey("When the reader reads a number of bytes", func() {
			Convey("Then it returns without error", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser with no unread bytes split into 1 reader", t, func() {
		Convey("When the reader reads a number of bytes", func() {
			Convey("Then it returns an EOF error", func() {})
		})
	})
}

func TestReadCloserSplit_Close(t *testing.T) {
	Convey("Given there is an upstream ReadCloser with unread bytes split into 1 reader", t, func() {
		Convey("When the reader is closed", func() {
			Convey("Then it closes without error", func() {})
			Convey("And the upstream ReadCloser is closed without error", func() {})
		})
	})

	Convey("Given there is an upstream ReadCloser with unread bytes split into 2 readers", t, func() {
		Convey("When the first reader is closed", func() {
			Convey("Then it closes without error", func() {})
			Convey("And the upstream ReadCloser is not closed", func() {})
		})

		Convey("When the second reader is closed", func() {
			Convey("Then it closes without error", func() {})
			Convey("And the upstream ReadCloser is closed without error", func() {})
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

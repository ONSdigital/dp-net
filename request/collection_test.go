package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testToken = "foo"
)

func TestGetCollectionID(t *testing.T) {
	Convey("should return collectionID from request header", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
		req.Header.Set(CollectionIDHeaderKey, testToken)

		actual, err := GetCollectionID(req)

		So(actual, ShouldEqual, testToken)
		So(err, ShouldBeNil)
	})

	Convey("should return collection id from request cookie", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
		req.AddCookie(&http.Cookie{Name: CollectionIDCookieKey, Value: testToken})

		actual, err := GetCollectionID(req)

		So(actual, ShouldEqual, testToken)
		So(err, ShouldBeNil)
	})

	Convey("should return empty token if no header or cookie is set", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)

		actual, err := GetCollectionID(req)

		So(actual, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("should return empty token and error if get header returns an error that is not ErrHeaderNotFound", t, func() {
		actual, err := GetCollectionID(nil)

		So(actual, ShouldBeEmpty)
		So(err, ShouldResemble, headers.ErrRequestNil)
	})

	Convey("should accept letters, digits and hyphens in a header", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
		req.Header.Set(CollectionIDHeaderKey, "AbC-123")

		actual, err := GetCollectionID(req)

		So(actual, ShouldEqual, "AbC-123")
		So(err, ShouldBeNil)
	})

	Convey("should accept letters, digits and hyphens in a cookie", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
		req.AddCookie(&http.Cookie{Name: CollectionIDCookieKey, Value: "AbC-123"})

		actual, err := GetCollectionID(req)

		So(actual, ShouldEqual, "AbC-123")
		So(err, ShouldBeNil)
	})

	for _, invalidID := range []string{"abc_def", "abc def", "abc.def", "abc/def", "caf?something"} {
		Convey("should reject invalid header collection id "+invalidID, t, func() {
			req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
			req.Header.Set(CollectionIDHeaderKey, invalidID)

			actual, err := GetCollectionID(req)

			So(actual, ShouldBeEmpty)
			So(err, ShouldEqual, ErrInvalidCollectionId)
		})

		Convey("should reject invalid cookie collection id "+invalidID, t, func() {
			req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
			req.AddCookie(&http.Cookie{Name: CollectionIDCookieKey, Value: invalidID})

			actual, err := GetCollectionID(req)

			So(actual, ShouldBeEmpty)
			So(err, ShouldEqual, ErrInvalidCollectionId)
		})
	}
}

func TestGetCollectionIDFromCookie(t *testing.T) {
	Convey("should return collection id from request cookie", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)
		req.AddCookie(&http.Cookie{Name: CollectionIDCookieKey, Value: testToken})

		actual, err := getCollectionIDFromCookie(req)

		So(actual, ShouldEqual, testToken)
		So(err, ShouldBeNil)
	})

	Convey("should return empty id if collection id cookie not found", t, func() {
		req := httptest.NewRequest("GET", "http://localhost:8080", http.NoBody)

		actual, err := getCollectionIDFromCookie(req)

		So(actual, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}

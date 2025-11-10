package fallback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testHeader           = "X-Test"
	primaryDesignation   = "primary"
	secondaryDesignation = "secondary"
	tertiaryDesignation  = "tertiary"
)

func TestAlternativeBuilderAndAlternativeServeHTTP(t *testing.T) {
	Convey("Given an AlternativeBuilder that tries a primary handler and uses a fallback handler on a trigger status", t, func() {
		primaryHandler := generateHandlerWithStatus(http.StatusNotFound, primaryDesignation)
		fallbackHandler := generateHandlerWithStatus(http.StatusOK, secondaryDesignation)

		alt := Try(primaryHandler).WhenStatus(http.StatusNotFound).Then(fallbackHandler)

		Convey("When the primary handler returns the trigger status", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", http.NoBody)
			alt.ServeHTTP(w, req)

			Convey("Then the fallback handler is called", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, secondaryDesignation+" response")
				So(w.Header().Get(testHeader), ShouldEqual, secondaryDesignation)
			})
		})

		Convey("When the primary handler returns a non-fallback status", func() {
			primaryOK := generateHandlerWithStatus(http.StatusOK, primaryDesignation)
			altOK := Try(primaryOK).WhenStatus(http.StatusNotFound).Then(fallbackHandler)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", http.NoBody)
			altOK.ServeHTTP(w, req)

			Convey("Then the primary handler's response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, primaryDesignation+" response")
				So(w.Header().Get(testHeader), ShouldEqual, primaryDesignation)
			})
		})
	})
}

func TestAlternativeBuilderChaining(t *testing.T) {
	Convey("Given a chained Alternative builder with three handlers", t, func() {
		primary := generateHandlerWithStatus(http.StatusNotFound, primaryDesignation)
		secondary := generateHandlerWithStatus(http.StatusForbidden, secondaryDesignation)
		final := generateHandlerWithStatus(http.StatusOK, tertiaryDesignation)

		chain := Try(primary).WhenStatus(http.StatusNotFound).Then(
			Try(secondary).WhenStatus(http.StatusForbidden).Then(final),
		)

		Convey("When the primary and secondary handlers return their respective fallback statuses", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", http.NoBody)
			chain.ServeHTTP(w, req)

			Convey("Then the final handler's response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, tertiaryDesignation+" response")
			})
		})
	})
}

func TestResponseWriter(t *testing.T) {
	Convey("Given a responseWriter implementation", t, func() {
		w := &responseWriter{}
		Convey("When setting headers and writing body", func() {
			w.Header().Set(testHeader, "value")
			w.WriteHeader(201)
			w.Write([]byte("body"))

			Convey("Then the header, status, and body are correct", func() {
				So(w.Header().Get(testHeader), ShouldEqual, "value")
				So(w.StatusCode(), ShouldEqual, 201)
				So(string(w.Body()), ShouldEqual, "body")
			})
		})
	})
}

func TestAlternativeWhenStatus(t *testing.T) {
	Convey("Given an Alternative, WhenStatus returns an AlternativeBuilder", t, func() {
		alt := &Alternative{}
		builder := alt.WhenStatus(http.StatusTeapot)
		So(builder, ShouldNotBeNil)
		So(*builder.whenStatus, ShouldEqual, http.StatusTeapot)
	})
}

func generateHandlerWithStatus(status int, designation string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(testHeader, designation)
		w.WriteHeader(status)
		w.Write([]byte(designation + " response"))
	})
}

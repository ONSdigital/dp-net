package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-net/handlers"
	. "github.com/smartystreets/goconvey/convey"
)

const testToken = "666"
const testLocale = "cy"

type mockHandler struct {
	invocations int
	ctx         context.Context
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.invocations++
	m.ctx = r.Context()
}

func TestCheckHeaderUserAccess(t *testing.T) {
	Convey("given the request with a florence access token header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.Header.Set(handlers.UserAccess.Header(), testToken)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckHeader(mockHandler, handlers.UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken, ok := mockHandler.ctx.Value(handlers.UserAccess.Context()).(string)
				So(ok, ShouldBeTrue)
				So(xFlorenceToken, ShouldEqual, testToken)
			})
		})
	})

	Convey("given the request does not contain a florence access token header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckHeader(mockHandler, handlers.UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken := mockHandler.ctx.Value(handlers.UserAccess.Context())
				So(xFlorenceToken, ShouldBeNil)
			})
		})
	})
}

func TestCheckHeaderLocale(t *testing.T) {
	Convey("given the request with a locale header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.Header.Set(handlers.Locale.Header(), testLocale)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckHeader(mockHandler, handlers.Locale)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key localeCode", func() {
				localeCode, ok := mockHandler.ctx.Value(handlers.Locale.Context()).(string)
				So(ok, ShouldBeTrue)
				So(localeCode, ShouldEqual, testLocale)
			})
		})
	})
}

func TestCheckCookieUserAccess(t *testing.T) {
	Convey("given the request contain a cookie for a florence access token header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.AddCookie(&http.Cookie{Name: handlers.UserAccess.Cookie(), Value: testToken})

		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckCookie(mockHandler, handlers.UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken, ok := mockHandler.ctx.Value(handlers.UserAccess.Context()).(string)
				So(ok, ShouldBeTrue)
				So(xFlorenceToken, ShouldEqual, testToken)
			})
		})
	})

	Convey("given the request does not contain a cookie for a florence access token header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckCookie(mockHandler, handlers.UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context does not contain value for key florence-id", func() {
				xFlorenceToken := mockHandler.ctx.Value(handlers.UserAccess.Context())
				So(xFlorenceToken, ShouldBeNil)
			})
		})
	})
}

func TestCheckCookieLocale(t *testing.T) {
	Convey("given the request contain a cookie for locale code ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.AddCookie(&http.Cookie{Name: handlers.Locale.Cookie(), Value: testLocale})

		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := handlers.DoCheckCookie(mockHandler, handlers.Locale)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key localeCode", func() {
				localeCode, ok := mockHandler.ctx.Value(handlers.Locale.Context()).(string)
				So(ok, ShouldBeTrue)
				So(localeCode, ShouldEqual, testLocale)
			})
		})
	})
}

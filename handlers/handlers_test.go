package handlers

import (
	"context"
	"fmt"
	dprequest "github.com/ONSdigital/dp-net/request"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testToken = "666"
const testLocale = "cy"
const testCollectionID = "foo"

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
		r.Header.Set(UserAccess.Header(), testToken)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := DoCheckHeader(mockHandler, UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken, ok := mockHandler.ctx.Value(UserAccess.Context()).(string)
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

		target := DoCheckHeader(mockHandler, UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken := mockHandler.ctx.Value(UserAccess.Context())
				So(xFlorenceToken, ShouldBeNil)
			})
		})
	})
}

func TestCheckHeaderLocale(t *testing.T) {
	Convey("given the request with a locale header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.Header.Set(Locale.Header(), testLocale)
		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := DoCheckHeader(mockHandler, Locale)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key localeCode", func() {
				localeCode, ok := mockHandler.ctx.Value(Locale.Context()).(string)
				So(ok, ShouldBeTrue)
				So(localeCode, ShouldEqual, testLocale)
			})
		})
	})
}

func TestCheckCookieUserAccess(t *testing.T) {
	Convey("given the request contain a cookie for a florence access token header ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.AddCookie(&http.Cookie{Name: UserAccess.Cookie(), Value: testToken})

		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := DoCheckCookie(mockHandler, UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key florence-id", func() {
				xFlorenceToken, ok := mockHandler.ctx.Value(UserAccess.Context()).(string)
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

		target := DoCheckCookie(mockHandler, UserAccess)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context does not contain value for key florence-id", func() {
				xFlorenceToken := mockHandler.ctx.Value(UserAccess.Context())
				So(xFlorenceToken, ShouldBeNil)
			})
		})
	})
}

func TestCheckCookieLocale(t *testing.T) {
	Convey("given the request contain a cookie for locale code ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.AddCookie(&http.Cookie{Name: Locale.Cookie(), Value: testLocale})

		w := httptest.NewRecorder()

		mockHandler := &mockHandler{
			invocations: 0,
		}

		target := DoCheckCookie(mockHandler, Locale)

		Convey("when the handler is called", func() {
			target.ServeHTTP(w, r)

			Convey("then the wrapped handle is called 1 time", func() {
				So(mockHandler.invocations, ShouldEqual, 1)
			})

			Convey("and the request context contains a value for key localeCode", func() {
				localeCode, ok := mockHandler.ctx.Value(Locale.Context()).(string)
				So(ok, ShouldBeTrue)
				So(localeCode, ShouldEqual, testLocale)
			})
		})
	})
}

func TestControllerHandler(t *testing.T) {
	Convey("given a controllerHandlerFunc and a context with a collectionID, florence ID and a cookie with a locale", t, func() {
		request := httptest.NewRequest("GET", "http://localhost:8080", nil)
		request = request.WithContext(context.WithValue(request.Context(), dprequest.CollectionIDHeaderKey, testCollectionID))
		request = request.WithContext(context.WithValue(request.Context(), dprequest.FlorenceIdentityKey, testToken))
		request.AddCookie(&http.Cookie{Name: Locale.Cookie(), Value: testLocale})
		w := httptest.NewRecorder()

		var controllerHandlerFunc ControllerHandlerFunc = func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
			r.Header.Set(UserAccess.Header(), accessToken)
			r.Header.Set(CollectionID.Header(), collectionID)
			r.Header.Set(Locale.Header(), lang)
		}

		h := ControllerHandler(controllerHandlerFunc)
		Convey("when the handler is called", func() {
			h.ServeHTTP(w, request)
			Convey("the request context contains a value for the locale code", func() {
				localeCode := request.Header.Get(Locale.Header())
				So(localeCode, ShouldEqual, testLocale)
			})
			Convey("the request context contains a value for the Collection identity", func() {
				localeCode := request.Header.Get(CollectionID.Header())
				So(localeCode, ShouldEqual, testCollectionID)
			})
			Convey("the request context contains a value for UserAccess token", func() {
				localeCode := request.Header.Get(UserAccess.Header())
				So(localeCode, ShouldEqual, testToken)
			})
		})
	})
	Convey("given a controllerHandlerFunc and a context with no collection ID or florence ID but with a locale in the subdomain", t, func() {
		target := fmt.Sprintf("http://%s.localhost:8080", testLocale)
		request := httptest.NewRequest("GET", target, nil)
		request.Header.Set(UserAccess.Header(), testToken)
		request.Header.Set(CollectionID.Header(), testCollectionID)
		w := httptest.NewRecorder()

		var controllerHandlerFunc ControllerHandlerFunc = func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
			r.Header.Set(UserAccess.Header(), accessToken)
			r.Header.Set(CollectionID.Header(), collectionID)
			r.Header.Set(Locale.Header(), lang)
		}

		h := ControllerHandler(controllerHandlerFunc)
		Convey("when the handler is called", func() {
			h.ServeHTTP(w, request)
			Convey("the request context contains a value for the locale code", func() {
				localeCode := request.Header.Get(Locale.Header())
				So(localeCode, ShouldEqual, testLocale)
			})
			Convey("the request context contains an empty string value for the Collection identity", func() {
				localeCode := request.Header.Get(CollectionID.Header())
				So(localeCode, ShouldEqual, "")
			})
			Convey("the request context contains an empty string value for UserAccess token", func() {
				localeCode := request.Header.Get(UserAccess.Header())
				So(localeCode, ShouldEqual, "")
			})
		})
	})

	Convey("given a controllerHandlerFunc with no context, no cookies and no subdomain", t, func() {
		request := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		var controllerHandlerFunc ControllerHandlerFunc = func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
			r.Header.Set(UserAccess.Header(), accessToken)
			r.Header.Set(CollectionID.Header(), collectionID)
			r.Header.Set(Locale.Header(), lang)
		}

		h := ControllerHandler(controllerHandlerFunc)
		Convey("when the handler is called", func() {
			h.ServeHTTP(w, request)

			testDesc := fmt.Sprintf("the request context contains a value for the default locale code: %s", dprequest.DefaultLang)
			Convey(testDesc, func() {
				localeCode := request.Header.Get(Locale.Header())
				So(localeCode, ShouldEqual, dprequest.DefaultLang)
			})
		})
	})
}

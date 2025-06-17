package http

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

// listenAndServeTLSCalls keeps track of a listenAndServeTLS call
type listenAndServeTLSCalls struct {
	httpServer *Server
	certFile   string
	keyFile    string
}

// listenAndServeCalls keeps track of a listenAndServe call
type listenAndServeCalls struct {
	httpServer *Server
}

func TestNew(t *testing.T) {

	Convey("Given mocked network calls", t, func() {

		doListenAndServe = func(httpServer *Server) error {
			return errors.New("unexpected ListenAndServe call")
		}

		doListenAndServeTLS = func(httpServer *Server, certFile, keyFile string) error {
			return errors.New("unexpected ListenAndServeTLS call")

		}

		doShutdown = func(ctx context.Context, httpServer *http.Server) error {
			return errors.New("unexpected Shutdown call")

		}

		Convey("New should return a new server with sensible defaults", func() {
			h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
			s := NewServer(":0", h)

			So(s, ShouldNotBeNil)
			So(s.Handler, ShouldEqual, h)
			So(s.Alice, ShouldBeNil)
			So(s.Addr, ShouldEqual, ":0")
			So(s.MaxHeaderBytes, ShouldEqual, 0)

			Convey("TLS should not be configured by default", func() {
				So(s.CertFile, ShouldBeEmpty)
				So(s.KeyFile, ShouldBeEmpty)
			})

			Convey("Default middleware should include RequestID and Log", func() {
				So(s.middleware, ShouldContainKey, RequestIDHandlerKey)
				So(s.middleware, ShouldContainKey, LogHandlerKey)
				So(s.middlewareOrder, ShouldResemble, []string{RequestIDHandlerKey, LogHandlerKey})
			})

			Convey("Default timeouts should be sensible", func() {
				So(s.ReadTimeout, ShouldEqual, time.Second*5)
				So(s.WriteTimeout, ShouldEqual, time.Second*10)
				So(s.ReadHeaderTimeout, ShouldEqual, 0)
				So(s.IdleTimeout, ShouldEqual, 0)
			})

			Convey("Handle OS signals by default", func() {
				So(s.HandleOSSignals, ShouldEqual, true)
			})

			Convey("A default shutdown context is initialised", func() {
				So(s.DefaultShutdownTimeout, ShouldEqual, 10*time.Second)
			})
		})

		Convey("prep should prepare the server correctly", func() {
			Convey("prep should create a valid Server instance", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)

				s.prep()
				So(s.Server.Addr, ShouldEqual, ":0")
			})

			Convey("invalid middleware should panic", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)

				s.middlewareOrder = []string{"foo"}

				So(func() {
					s.prep()
				}, ShouldPanicWith, "middleware not found: foo")
			})

			Convey("ListenAndServe with invalid middleware should panic", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)

				s.middlewareOrder = []string{"foo"}
				var err error
				So(func() {
					err = s.ListenAndServe()
				}, ShouldPanicWith, "middleware not found: foo")

				So(err, ShouldBeNil)
			})

			Convey("ListenAndServeTLS with invalid middleware should panic", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)

				s.middlewareOrder = []string{"foo"}
				var err error
				So(func() {
					err = s.ListenAndServeTLS("testdata/certFile", "testdata/keyFile")
				}, ShouldPanicWith, "middleware not found: foo")

				So(err, ShouldBeNil)
			})
		})

		Convey("ListenAndServeTLS", func() {
			Convey("ListenAndServeTLS should set CertFile/KeyFile", func() {
				wg := &sync.WaitGroup{}
				calls := []listenAndServeTLSCalls{}
				doListenAndServeTLS = func(httpServer *Server, certFile, keyFile string) error {
					defer wg.Done()
					calls = append(calls, listenAndServeTLSCalls{
						httpServer: httpServer,
						certFile:   certFile,
						keyFile:    keyFile})
					return nil
				}

				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)
				var err error
				// execute ListenAndServer and wait for it to finish
				wg.Add(1)
				go func() {
					err = s.ListenAndServeTLS("testdata/certFile", "testdata/keyFile")
				}()
				wg.Wait()

				So(s.CertFile, ShouldEqual, "testdata/certFile")
				So(s.KeyFile, ShouldEqual, "testdata/keyFile")
				So(calls, ShouldHaveLength, 1)
				So(calls[0].httpServer, ShouldNotBeNil)
				So(calls[0].certFile, ShouldEqual, "testdata/certFile")
				So(calls[0].keyFile, ShouldEqual, "testdata/keyFile")

				So(err, ShouldBeNil)
			})

			Convey("ListenAndServeTLS with only CertFile should panic", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)
				var err error
				So(func() {
					err = s.ListenAndServeTLS("certFile", "")
				}, ShouldPanicWith, "either CertFile/KeyFile must be blank, or both provided")

				So(err, ShouldBeNil)
			})

			Convey("ListenAndServeTLS with only KeyFile should panic", func() {
				h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
				s := NewServer(":0", h)
				var err error
				So(func() {
					err = s.ListenAndServeTLS("", "keyFile")
				}, ShouldPanicWith, "either CertFile/KeyFile must be blank, or both provided")

				So(err, ShouldBeNil)
			})
		})

		Convey("Given a mocked ListenAndServe", func() {
			wg := &sync.WaitGroup{}
			calls := []listenAndServeCalls{}
			doListenAndServe = func(httpServer *Server) error {
				defer wg.Done()
				calls = append(calls, listenAndServeCalls{httpServer: httpServer})
				return nil
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
			s := NewServer(":0", h)

			Convey("then ListenAndServe starts a working HTTP server", func() {
				So(s.HandleOSSignals, ShouldBeTrue)
				var err error
				wg.Add(1)
				go func() {
					err = s.ListenAndServe()
				}()
				wg.Wait()

				So(calls, ShouldHaveLength, 1)
				So(calls[0].httpServer, ShouldNotBeNil)

				So(err, ShouldBeNil)
			})

			Convey("then if HandleOSSignals is disabled, ListenAndServe starts a working HTTP server", func() {
				s.HandleOSSignals = false
				var err error
				wg.Add(1)
				go func() {
					err = s.ListenAndServe()
				}()
				wg.Wait()

				So(calls, ShouldHaveLength, 1)
				So(calls[0].httpServer, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestServer_LongRunningOperation(t *testing.T) {
	doListenAndServe = func(httpServer *Server) error {
		return timeoutHandler(httpServer).ListenAndServe()
	}
	doShutdown = func(ctx context.Context, httpServer *http.Server) error {
		return httpServer.Shutdown(ctx)
	}

	writeTimeout := 1 * time.Second
	Convey("given a free port on the localhost", t, func() {
		p, err := GetFreePort()
		if err != nil {
			t.Fatalf("Cannot find a free port to perform test: %v", err)
		}
		a := "localhost:" + strconv.Itoa(p)

		Convey("with a server whose handler runs for less than the server's WriteTimout", func() {
			h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				_, _ = w.Write([]byte("Done"))
			})
			eChan, cleanup := startServer(a, h, writeTimeout, 0)
			select {
			case err = <-eChan:
				So(err, ShouldBeNil)
			case <-time.After(10 * time.Millisecond):
				defer cleanup()
			}

			Convey("when a request is made to the server", func() {
				resp, err := http.Get("http://" + a)

				Convey("the request is completed successfully", func() {
					So(err, ShouldEqual, nil)
					So(resp.StatusCode, ShouldEqual, 200)

					b, err := io.ReadAll(resp.Body)
					So(err, ShouldBeNil)
					So(string(b), ShouldEqual, "Done")
				})
			})
		})

		Convey("with a server whose handler runs for longer than the server's WriteTimout (no response timeout)", func() {
			h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				time.Sleep(writeTimeout + 100*time.Millisecond)
				_, _ = w.Write([]byte("Done"))
			})
			eChan, cleanup := startServer(a, h, writeTimeout, 0)
			select {
			case err = <-eChan:
				So(err, ShouldBeNil)
			case <-time.After(10 * time.Millisecond):
				defer cleanup()
			}

			Convey("when a request is made to the server", func() {
				resp, err := http.Get("http://" + a)

				Convey("the request is terminated", func() {
					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
			})
		})

		Convey("with a server whose handler runs for longer than the server's WriteTimout (with response timeout)", func() {
			h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				time.Sleep(writeTimeout + 100*time.Millisecond)
				_, _ = w.Write([]byte("Done"))
			})
			// this will test that there will always be sufficient time
			// to write to reponse
			eChan, cleanup := startServer(a, h, writeTimeout, writeTimeout)
			select {
			case err = <-eChan:
				So(err, ShouldBeNil)
			case <-time.After(10 * time.Millisecond):
				defer cleanup()
			}

			Convey("when a request is made to the server", func() {
				resp, err := http.Get("http://" + a)

				Convey("the request is terminated with a 'test server timeout' response", func() {
					So(err, ShouldBeNil)
					So(resp.StatusCode, ShouldEqual, http.StatusServiceUnavailable)

					b, err := io.ReadAll(resp.Body)
					So(err, ShouldEqual, nil)
					So(string(b), ShouldEqual, "test server timeout")
				})
			})
		})
	})
}

func TestGetFreePort(t *testing.T) {
	Convey("When GetFreePort() is called n times, where n > 1", t, func() {
		n := 10

		Convey("A free, usable port should be returned every time", func() {
			for i := 0; i < n; i++ {
				port, err := GetFreePort()
				So(err, ShouldBeNil)
				So(port, ShouldNotEqual, 0)

				l, e := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
				So(e, ShouldBeNil)
				So(l, ShouldNotBeNil)
				_ = l.Close()
			}
		})
	})
}

func startServer(address string, handler http.Handler, writeTimeout time.Duration, requestTimeout time.Duration) (chan error, func()) {
	var s *Server
	if requestTimeout > 0 {
		s = NewServerWithTimeout(address, handler, requestTimeout, "test server timeout")
	} else {
		s = NewServer(address, handler)
	}
	s.WriteTimeout = writeTimeout
	s.RequestTimeout = requestTimeout
	s.HandleOSSignals = false

	eChan := make(chan error)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			eChan <- err
		}
		close(eChan)
	}()

	return eChan, func() {
		if err := s.Shutdown(context.Background()); err != nil {
			So(err, ShouldBeNil)
		}
	}
}

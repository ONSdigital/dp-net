package http

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	request "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/justinas/alice"
)

const (
	RequestIDHandlerKey string = "RequestID"
	LogHandlerKey       string = "Log"
)

var (
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

// Server is a http.Server with sensible defaults, which supports
// configurable middleware and timeouts, and shuts down cleanly
// on SIGINT/SIGTERM
type Server struct {
	http.Server
	middleware             map[string]alice.Constructor
	middlewareOrder        []string
	Alice                  *alice.Chain
	CertFile               string
	KeyFile                string
	DefaultShutdownTimeout time.Duration
	HandleOSSignals        bool
}

// NewServer creates a new server
func NewServer(bindAddr string, router http.Handler) *Server {
	middleware := map[string]alice.Constructor{
		RequestIDHandlerKey: request.HandlerRequestID(16),
		LogHandlerKey:       log.Middleware,
	}

	return &Server{
		Alice:           nil,
		middleware:      middleware,
		middlewareOrder: []string{RequestIDHandlerKey, LogHandlerKey},
		Server: http.Server{
			Addr:         bindAddr,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
			// Give the timeout handler some time to write the response before the writer is closed
			Handler:           http.TimeoutHandler(router, DefaultWriteTimeout-100*time.Millisecond, "connection timeout"),
			ReadHeaderTimeout: 0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
		},
		HandleOSSignals:        true,
		DefaultShutdownTimeout: 10 * time.Second,
	}
}

func (s *Server) prep() {
	var m []alice.Constructor
	for _, v := range s.middlewareOrder {
		if mw, ok := s.middleware[v]; ok {
			m = append(m, mw)
			continue
		}
		panic("middleware not found: " + v)
	}

	s.Server.Handler = alice.New(m...).Then(s.Handler)
}

// ListenAndServe sets up SIGINT/SIGTERM signals, builds the middleware
// chain, and creates/starts a http.Server instance
//
// If CertFile/KeyFile are both set, the http.Server instance is started
// using ListenAndServeTLS. Otherwise ListenAndServe is used.
//
// Specifying one of CertFile/KeyFile without the other will panic.
func (s *Server) ListenAndServe() error {
	if s.HandleOSSignals {
		return s.listenAndServeHandleOSSignals()
	}

	return s.listenAndServe()
}

// ListenAndServeTLS sets KeyFile and CertFile, then calls ListenAndServe
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if len(certFile) == 0 || len(keyFile) == 0 {
		panic("either CertFile/KeyFile must be blank, or both provided")
	}
	s.KeyFile = keyFile
	s.CertFile = certFile
	return s.ListenAndServe()
}

// Shutdown will gracefully shutdown the server, using a default shutdown
// timeout if a context is not provided.
func (s *Server) Shutdown(ctx context.Context) error {

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), s.DefaultShutdownTimeout)
	}

	return doShutdown(ctx, &s.Server)
}

func (s *Server) listenAndServe() error {

	s.prep()
	if len(s.CertFile) > 0 || len(s.KeyFile) > 0 {
		return doListenAndServeTLS(&s.Server, s.CertFile, s.KeyFile)
	}

	return doListenAndServe(&s.Server)
}

func (s *Server) listenAndServeHandleOSSignals() error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	s.listenAndServeAsync()

	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return s.Shutdown(ctx)
}

func (s *Server) listenAndServeAsync() {

	s.prep()
	if len(s.CertFile) > 0 || len(s.KeyFile) > 0 {
		go func() {
			if err := doListenAndServeTLS(&s.Server, s.CertFile, s.KeyFile); err != nil {
				log.Error(context.Background(), "http server returned error", err)
				os.Exit(1)
			}
		}()
	} else {
		go func() {
			if err := doListenAndServe(&s.Server); err != nil {
				log.Error(context.Background(), "http server returned error", err)
				os.Exit(1)
			}
		}()
	}
}

var doListenAndServe = func(httpServer *http.Server) error {
	return httpServer.ListenAndServe()
}

var doListenAndServeTLS = func(httpServer *http.Server, certFile, keyFile string) error {
	return httpServer.ListenAndServeTLS(certFile, keyFile)
}

var doShutdown = func(ctx context.Context, httpServer *http.Server) error {
	return httpServer.Shutdown(ctx)
}

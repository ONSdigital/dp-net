package http

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	"github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/justinas/alice"
)

const (
	RequestIDHandlerKey string        = "RequestID"
	LogHandlerKey       string        = "Log"
	ResponseWriteGrace  time.Duration = 100 * time.Millisecond
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
	RequestTimeout         time.Duration
	TimeoutMessage         string
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
			Handler:           router,
			Addr:              bindAddr,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
		},
		HandleOSSignals:        true,
		DefaultShutdownTimeout: 10 * time.Second,
	}
}

// NewServerWithTimeout creates a new server with request timeout duration
// and a message that will be in the response body
func NewServerWithTimeout(bindAddr string, router http.Handler, timeout time.Duration, timeoutMessage string) *Server {
	server := NewServer(bindAddr, router)
	server.RequestTimeout = timeout
	server.TimeoutMessage = timeoutMessage
	return server
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

	s.Handler = alice.New(m...).Then(s.Handler)
}

// ListenAndServe sets up SIGINT/SIGTERM signals, builds the middleware
// chain, and creates/starts a http.Server instance
//
// If CertFile/KeyFile are both set, the http.Server instance is started
// using ListenAndServeTLS. Otherwise, ListenAndServe is used.
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
	if certFile == "" || keyFile == "" {
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
	if s.CertFile != "" || s.KeyFile != "" {
		return doListenAndServeTLS(s, s.CertFile, s.KeyFile)
	}

	return doListenAndServe(s)
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
	if s.CertFile != "" || s.KeyFile != "" {
		go func() {
			if err := doListenAndServeTLS(s, s.CertFile, s.KeyFile); err != nil {
				log.Error(context.Background(), "http server returned error", err)
				os.Exit(1)
			}
		}()
	} else {
		go func() {
			if err := doListenAndServe(s); err != nil {
				log.Error(context.Background(), "http server returned error", err)
				os.Exit(1)
			}
		}()
	}
}

func timeoutHandler(s *Server) *http.Server {
	if s.RequestTimeout > 0 {
		// give some time for the response to be written
		if s.WriteTimeout <= s.RequestTimeout {
			s.WriteTimeout = s.RequestTimeout + ResponseWriteGrace
		}
		timeoutMsg := "connection timeout"
		if s.TimeoutMessage != "" {
			timeoutMsg = s.TimeoutMessage
		}
		s.Handler = http.TimeoutHandler(s.Handler, s.RequestTimeout, timeoutMsg)
	}

	return &s.Server
}

var doListenAndServe = func(httpServer *Server) error {
	return timeoutHandler(httpServer).ListenAndServe()
}

var doListenAndServeTLS = func(httpServer *Server, certFile, keyFile string) error {
	return timeoutHandler(httpServer).ListenAndServeTLS(certFile, keyFile)
}

var doShutdown = func(ctx context.Context, httpServer *http.Server) error {
	return httpServer.Shutdown(ctx)
}

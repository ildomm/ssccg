package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	DefaultListenAddress     = 8080
	DefaultReadHeaderTimeout = time.Second * 15
	DefaultWriteTimeout      = time.Second * 15
	DefaultReadTimeout       = time.Second * 15
	DefaultIdleTimeout       = time.Second * 60
)

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress     int
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	readTimeout       time.Duration
	idleTimeout       time.Duration
}

// NewServer is a factory to instantiate a new Server.
func NewServer() *Server {

	return &Server{
		listenAddress:     DefaultListenAddress,
		readHeaderTimeout: DefaultReadHeaderTimeout,
		writeTimeout:      DefaultWriteTimeout,
		readTimeout:       DefaultReadTimeout,
		idleTimeout:       DefaultIdleTimeout,
	}
}

// Run defines the server and starts it.
func (s *Server) Run() error {

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", s.listenAddress),

		// Good practice to set timeouts to avoid Slow-loris attacks.
		ReadHeaderTimeout: s.readHeaderTimeout,
		WriteTimeout:      s.writeTimeout,
		ReadTimeout:       s.readTimeout,
		IdleTimeout:       s.idleTimeout,

		Handler: s.router(),
	}

	return httpServer.ListenAndServe()
}

// router registers all HandlerFunc and middleware for the existing HTTP routes.
func (s *Server) router() *mux.Router {

	rootRouter := mux.NewRouter()
	rootRouter.HandleFunc("/api/v0/health", s.HealthHandler)

	rootRouter.Use(NewRecoverMiddleware())
	rootRouter.Use(NewLoggingMiddleware())

	return rootRouter
}

func (s *Server) ListenAddress() int {
	return s.listenAddress
}

func (s *Server) WithListenAddress(listenAddress int) {
	s.listenAddress = listenAddress
}

func (s *Server) WithReadHeaderTimeout(readHeaderTimeout time.Duration) {
	s.readHeaderTimeout = readHeaderTimeout
}

func (s *Server) WithWriteTimeout(writeTimeout time.Duration) {
	s.writeTimeout = writeTimeout
}

func (s *Server) WithReadTimeout(readTimeout time.Duration) {
	s.readTimeout = readTimeout
}

func (s *Server) WithIdleTimeout(idleTimeout time.Duration) {
	s.idleTimeout = idleTimeout
}

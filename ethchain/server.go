package ethchain

import (
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	// points the the router handlers
	handler http.Handler
	// Listens for HTTP traffic on this address
	listenAddress string

	// http server
	srv *http.Server
}

func (s *Server) Initialize(host string, port uint16, allowedOrigins []string, handler http.Handler) {
	s.listenAddress = fmt.Sprintf("%s:%d", host, port)
	s.handler = handler
}

// Dispatch starts the API server
func (s *Server) Dispatch() error {
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	s.srv = &http.Server{Handler: s.handler}
	return s.srv.Serve(listener)
}

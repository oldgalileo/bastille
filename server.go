package main

import (
	"net"

	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	return &Server{
		http.NewServeMux(),
	}
}

func (s *Server) init() {
	listener, listenerErr := net.Listen("tcp", flagServerPort)
	if listenerErr != nil {
		log.Panic(listenerErr)
	}
	for {
		conn, connErr := listener.Accept()
		if connErr != nil {
			log.Error("Error accepting incoming request")
		}

	}
}

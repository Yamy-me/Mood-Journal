package moodJournal

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) (Server, error) {
	server := &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
	}
	return Server{httpServer: server}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

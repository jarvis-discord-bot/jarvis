package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/pedromsmoreira/jarvis/internal/configuration"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type server struct {
	settings *configuration.Settings
	router   *mux.Router
	server   *http.Server
}

func NewServer(settings *configuration.Settings, router *mux.Router) *server {
	return &server{
		settings: settings,
		router:   router,
	}
}

func (s *server) Start(wg *sync.WaitGroup) (string, error) {
	addr := fmt.Sprintf("%s:%d", s.settings.Srv.Address, s.settings.Srv.Port)
	s.server = &http.Server{
		Handler: s.router,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return addr, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.server.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Error("HTTP Server closed unexpectedly.")
		}
	}()
	return addr, nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

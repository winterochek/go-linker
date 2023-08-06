package server

import (
	"context"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/winterochek/go-linker/internal/shorten"
)

type CloseFunc func(context.Context) error

type Server struct {
	e         *echo.Echo
	shortener *shorten.Service
	closers   []CloseFunc
}

func New(shortener *shorten.Service) *Server {
	s := &Server{
		shortener: shortener,
	}
	s.setupRouter()

	return s
}

func (s *Server) AddCloser(closer CloseFunc) {
	s.closers = append(s.closers, closer)
}

func (s *Server) setupRouter() {
	s.e = echo.New()
	s.e.HideBanner = true
	s.e.Validator = NewValidator()

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

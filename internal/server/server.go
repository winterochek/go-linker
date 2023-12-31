package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/winterochek/go-linker/internal/config"
	"github.com/winterochek/go-linker/internal/model"
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

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Pre(middleware.RequestID())

	restricted := s.e.Group("/api")
	{
		restricted.Use(middleware.JWTWithConfig(makeJWTConfig()))
		restricted.POST("/shorten", HandleShorten(s.shortener))
	}

	s.e.GET("/:identifier", HandleRedirect(s.shortener))

	s.AddCloser(s.e.Shutdown)

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

func makeJWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &model.UserClaims{},
		SigningKey: []byte(config.Get().Auth.JWTSecretKey),
		ErrorHandler: func(err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized)
		},
	}
}

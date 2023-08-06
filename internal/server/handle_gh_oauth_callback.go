package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winterochek/go-linker/internal/config"
	"github.com/winterochek/go-linker/internal/model"
)

type callbackProvider interface {
	GitHubAuthCallback(ctx context.Context, sessionCode string) (*model.User, string, error)
}

func HandleGitHubAuthCallback(cbProvider callbackProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionCode := c.QueryParam("code")
		if sessionCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missing session code")
		}

		_, jwt, err := cbProvider.GitHubAuthCallback(c.Request().Context(), sessionCode)
		if err != nil {
			log.Printf("error handling github auth callback: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		redirectURL := fmt.Sprintf("%s/auth/token.html?token=%s", config.Get().BaseURL, jwt)
		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}

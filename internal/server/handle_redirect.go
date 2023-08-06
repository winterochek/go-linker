package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/winterochek/go-linker/internal/model"
)

type redirecter interface {
	Redirect(ctx context.Context, identifier string) (string, error)
}

func HandleRedirect(redirecter redirecter) echo.HandlerFunc {
	return func(c echo.Context) error {
		identifier := c.Param("identifier")

		redirectURL, err := redirecter.Redirect(c.Request().Context(), identifier)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			log.Printf("failed to redirect: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}

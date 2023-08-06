package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	. "github.com/samber/mo"
	"github.com/winterochek/go-linker/internal/config"
	"github.com/winterochek/go-linker/internal/model"
	"github.com/winterochek/go-linker/internal/shorten"
)

type shortener interface {
	Shorten(context.Context, model.ShortenInput) (*model.Shortening, error)
}

type shortenRequest struct {
	URL        string `json:"url" validate:"required,url"`
	Identifier string `json:"identifier,omitempty" validate:"omitempty,alphanum"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url,omitempty"`
	Message  string `json:"message,omitempty"`
}

func HandleShorten(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req shortenRequest

		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		identifier := None[string]()
		if strings.TrimSpace(req.Identifier) != "" {
			identifier = Some(req.Identifier)
		}

		input := model.ShortenInput{
			RawURL:     req.URL,
			Identifier: identifier,
		}

		shortening, err := shortener.Shorten(c.Request().Context(), input)
		if err != nil {
			if errors.Is(err, model.ErrIdentifierExists) {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}

			log.Printf("shortener.Shorten: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		shortURL, err := shorten.PrependBaseUrl(config.Get().BaseURL, shortening.Identifier)
		if err != nil {
			log.Printf("shorten.PrependBaseURL: $v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, shortenResponse{ShortURL: shortURL})
	}
}

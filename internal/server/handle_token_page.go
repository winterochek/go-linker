package server

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v48/github"
	"github.com/labstack/echo/v4"
	"github.com/winterochek/go-linker/internal/config"
	"github.com/winterochek/go-linker/internal/model"
)

// go:embed static
var static embed.FS

func HandleTokenPage() echo.HandlerFunc {
	tmpl, err := template.ParseFS(static, "static/token-page.html")
	if err != nil {
		log.Fatalf("failed to parse token.html template: %v", err)
	}

	type templateData struct {
		Token                   string
		TelegramContactUsername string
		GitHubUsername          string
		GitHubAvatarURL         string
	}

	type request struct {
		Token string `query:"token"`
	}
	return func(c echo.Context) error {
		var req request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		tokenObj, err := jwt.ParseWithClaims(req.Token, &model.UserClaims{}, func(_ *jwt.Token) (any, error) { return []byte(config.Get().Auth.JWTSecretKey), nil })
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		if !tokenObj.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := tokenObj.Claims.(*model.UserClaims)
		if !ok {
			log.Printf("failed to cast claims to *jwt.UserClaims")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		ghClient := github.NewClient(nil)

		ghUser, _, err := ghClient.Users.Get(c.Request().Context(), claims.User.GitHubLogin)
		if err != nil {
			log.Printf("failed to get github user: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		var (
			buf  bytes.Buffer
			data = templateData{
				Token:                   req.Token,
				TelegramContactUsername: config.Get().TelegramContactUsername,
				GitHubUsername:          ghUser.GetLogin(),
				GitHubAvatarURL:         *ghUser.AvatarURL,
			}
		)

		if err := tmpl.Execute(&buf, data); err != nil {
			log.Printf("failed to execute token.html template: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.HTML(http.StatusOK, buf.String())

	}
}

func HandleStatic() echo.HandlerFunc {
	return echo.WrapHandler(http.FileServer(http.FS(static)))
}

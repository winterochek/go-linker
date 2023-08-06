package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/winterochek/go-linker/internal/config"
	"github.com/winterochek/go-linker/internal/model"
)

func MakeJWT(user model.User) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "defer panic",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().Auth.JWTSecretKey))
}
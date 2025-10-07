package jwtToken

import (
	"errors"
	"time"

	"github.com/goggle-source/grpc-servic/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func GetToken(user domain.User, app domain.App, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(exp).Unix()
	claims["app_id"] = app.ID

	if app.Secret == "" {
		return "", errors.New("error secretKey")
	}

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

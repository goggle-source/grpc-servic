package jwtToken

import (
	"testing"
	"time"

	"github.com/goggle-source/grpc-servic/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func TestJwtTo(t *testing.T) {
	type test struct {
		name  string
		user  domain.User
		app   domain.App
		exp   time.Duration
		IsErr bool
	}

	tests := []test{
		{
			name: "success",
			user: domain.User{
				Email: "jonn@gmail.com",
			},
			app: domain.App{
				ID:     2,
				Secret: "tokenSecret",
			},
			exp:   10 * time.Second,
			IsErr: false,
		},
		{
			name: "error secretKey",
			user: domain.User{
				Email: "jonn@gmail.com",
			},
			app: domain.App{
				ID: 12,
			},
			exp:   10 * time.Minute,
			IsErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			token, err := GetToken(test.user, test.app, test.exp)
			if test.IsErr {
				if err == nil {
					t.Fatal("field is error")
				}
			} else {

				if err != nil {
					t.Errorf("invalid get token")
				}
			}

			if token != "" {
				
				JWTtoken, err := ParseToken(token, test.app.Secret)
				if err != nil {
					t.Fatal("field parse token")
				}

				if JWTtoken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					t.Error("invalid method")
				}
			}

		})
	}
}

func ParseToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
}

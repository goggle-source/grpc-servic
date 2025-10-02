package auth

import (
	"github.com/goggle-source/grpc-servic/sso/internal/grpc/auth"
)

type ServerAPI struct {
	auth.UnimplementedAuthServer
}

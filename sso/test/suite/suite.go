package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/goggle-source/grpc-servic/protos/gen/go/sso"
	"github.com/goggle-source/grpc-servic/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcHost = "localhost"

type Suilte struct {
	*testing.T
	Cfg        config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suilte) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(*cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection error: %v", err)
	}

	return ctx, &Suilte{
		Cfg:        *cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

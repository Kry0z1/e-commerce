package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	ssov1 "github.com/Kry0z1/e-commerce/protos/gen/go/sso"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Auth ssov1.AuthClient
	Cfg  *config.Config
}

func New(t *testing.T) (context.Context, Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(ctx, grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect to grpc server: %v", err)
	}

	return ctx, Suite{
		T:    t,
		Auth: ssov1.NewAuthClient(cc),
		Cfg:  cfg,
	}
}

func configPath() string {
	var res string
	if res = os.Getenv("CONFIG_PATH"); res == "" {
		res = "../config/local_tests.yaml"
	}

	return res
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port))
}

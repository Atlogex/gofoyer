package suite

import (
	ssov1 "atlogex/gofoyer/contractor/gen/go/sso"
	"atlogex/gofoyer/internal/config"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadBypath("/../config/local_test.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPCTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	address := grpcAddress(cfg)
	t.Logf("grpc address: %s", address)

	//cc, err := grpc.NewClient(
	//    address,
	//    grpc.WithTransportCredentials(insecure.NewCredentials()))
	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

const grpcHost = "localhost"

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GPRC.Port))
}

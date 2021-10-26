package grpc

import (
	"context"
	"fmt"
	"os"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	authority "github.com/dai65527/microservice-handson/services/authority/proto"
	catalog "github.com/dai65527/microservice-handson/services/catalog/proto"
	"github.com/dai65527/microservice-handson/services/gateway/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	ahost := os.Getenv("AUTHORITY_HOST")
	if ahost == "" {
		ahost = "localhost"
	}

	aport := os.Getenv("AUTHORITY_PORT")
	if aport == "" {
		aport = "5003"
	}

	aconn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", ahost, aport), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial authority grpc server: %w", err)
	}

	chost := os.Getenv("CUSTOMER_HOST")
	if chost == "" {
		chost = "localhost"
	}

	cport := os.Getenv("CUSTOMER_PORT")
	if aport == "" {
		aport = "5004"
	}

	cconn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", chost, cport), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial catalog grpc server: %w", err)
	}

	svc := &server{
		authorityClient: authority.NewAuthorityServiceClient(aconn),
		catalogClient:   catalog.NewCatalogServiceClient(cconn),
		logger:          logger.WithName("server"),
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterGatewayServiceServer(s, svc)
	}).Start(ctx)
}

package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	"github.com/dai65527/microservice-handson/services/authority/proto"
	customer "github.com/dai65527/microservice-handson/services/customer/proto"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	customerHost := os.Getenv("CUSTOMER_HOST")
	if customerHost == "" {
		customerHost = "localhost"
	}
	customerPort := os.Getenv("CUSTOMER_PORT")
	if customerPort == "" {
		customerPort = "5002"
	}
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", customerHost, customerPort), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial grpc customer server: %w", err)
	}

	customerClient := customer.NewCustomerServiceClient(conn)

	svc := &server{
		customerClient: customerClient,
		logger:         logger.WithName("server"),
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterAuthorityServiceServer(s, svc)
	}).Start(ctx)
}

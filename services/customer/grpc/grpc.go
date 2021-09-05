package grpc

import (
	"context"
	"fmt"
	"os"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	db "github.com/dai65527/microservice-handson/platform/db/proto"
	"github.com/dai65527/microservice-handson/services/customer/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5000"
	}
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", dbHost, dbPort), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial db server: %w", err)
	}

	dbClient := db.NewDBServiceClient(conn)
	svc := &server{
		dbClient: dbClient,
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterCustomerServiceServer(s, svc)
	}).Start(ctx)
}

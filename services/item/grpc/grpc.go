package grpc

import (
	"context"
	"fmt"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	db "github.com/dai65527/microservice-handson/platform/db/proto"
	"github.com/dai65527/microservice-handson/services/item/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	// conn, err := grpc.DialContext(ctx, "db.db.svc.cluster.local:5000", opts...)
	conn, err := grpc.DialContext(ctx, "localhost:5000", opts...)
	if err != nil {
		return fmt.Errorf("failed to dial db server: %w", err)
	}

	dbClinent := db.NewDBServiceClient(conn)
	svc := &server{
		dbClient: dbClinent,
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterItemServiceServer(s, svc)
	}).Start(ctx)
}

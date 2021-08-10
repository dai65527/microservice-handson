package grpc

import (
	"context"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	"github.com/dai65527/microservice-handson/platform/db/db"
	"github.com/dai65527/microservice-handson/platform/db/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	svc := &server{
		db: db.New(),
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterDBServiceServer(s, svc)
	}).Start(ctx)
}

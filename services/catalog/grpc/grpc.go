package grpc

import (
	"context"
	"fmt"
	"os"

	pkggrpc "github.com/dai65527/microservice-handson/pkg/grpc"
	"github.com/dai65527/microservice-handson/services/catalog/proto"
	customer "github.com/dai65527/microservice-handson/services/customer/proto"
	item "github.com/dai65527/microservice-handson/services/item/proto"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger logr.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	itemHost := os.Getenv("ITEM_HOST")
	if itemHost == "" {
		itemHost = "localhost"
	}
	itemPort := os.Getenv("ITEM_PORT")
	if itemPort == "" {
		itemPort = "5001"
	}
	// item connection
	iconn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", itemHost, itemPort), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial item grpc server: %w", err)
	}

	customerHost := os.Getenv("CUSTOMER_HOST")
	if customerHost == "" {
		customerHost = "localhost"
	}
	customerPort := os.Getenv("CUSTOMER_PORT")
	if customerPort == "" {
		customerPort = "5002"
	}
	// customer connection
	cconn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", customerHost, customerPort), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial customer grpc server: %w", err)
	}

	svc := &server{
		customerClient: customer.NewCustomerServiceClient(cconn),
		itemClient:     item.NewItemServiceClient(iconn),
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterCatalogServiceServer(s, svc)
	}).Start(ctx)
}

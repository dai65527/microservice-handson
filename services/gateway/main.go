package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/110y/run"
	"github.com/dai65527/microservice-handson/pkg/logger"
	"github.com/dai65527/microservice-handson/services/gateway/grpc"
	"github.com/dai65527/microservice-handson/services/gateway/http"
)

func main() {
	run.Run(server)
}

func server(ctx context.Context) int {
	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		if ferr != nil {
			// Unhandleable, something went wrong...
			panic(fmt.Sprintf("failed to write log:`%s` original error is:`%s`", ferr, err))
		}
		return 1
	}
	glogger := l.WithName("gateway")

	grpcPortStr := os.Getenv("GATEWAYGRPC_PORT")
	if grpcPortStr == "" {
		grpcPortStr = "5005"
	}
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse GATEWAYGRPC_PORT: %s", err))
	}

	httpPortStr := os.Getenv("GATEWAYHTTP_PORT")
	if httpPortStr == "" {
		httpPortStr = "4000"
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse GATEWAYHTTP_PORT: %s", err))
	}

	grpcErrCh := make(chan error, 1)
	go func() {
		grpcErrCh <- grpc.RunServer(ctx, grpcPort, glogger.WithName("grpc"))
	}()

	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- http.RunServer(ctx, httpPort, grpcPort)
	}()

	select {
	case err := <-grpcErrCh:
		glogger.Error(err, "failed to serve grpc server")
		return 1
	case err := <-httpErrCh:
		glogger.Error(err, "failed to serve http server")
		return 1
	case <-ctx.Done():
		glogger.Info("shutting down...")
		return 0
	}
}

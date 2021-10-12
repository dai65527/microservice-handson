package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/110y/run"
	"github.com/dai65527/microservice-handson/pkg/logger"
	"github.com/dai65527/microservice-handson/services/item/grpc"
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

	grpcErrCh := make(chan error, 1)
	go func() {
		portStr := os.Getenv("GATEWAY_PORT")
		if portStr == "" {
			portStr = "5005"
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			grpcErrCh <- fmt.Errorf("failed to parse GATEWAY_PORT: %s", err)
		}
		grpcErrCh <- grpc.RunServer(ctx, port, glogger.WithName("grpc"))
	}()

	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- http.RunServer(ctx, )
	}

	return 0
}

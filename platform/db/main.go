package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/dai65527/microservice-handson/pkg/logger"
	"github.com/dai65527/microservice-handson/platform/db/grpc"
	"golang.org/x/sys/unix"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
	defer stop()

	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		if ferr != nil {
			// Unhandleable, something went wrong...
			panic(fmt.Sprintf("failed to write log:`%s` original error is:`%s`", ferr, err))
		}
		return 1
	}
	clogger := l.WithName("db")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to parse port: %s", err))
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpc.RunServer(ctx, portInt, clogger.WithName("grpc"))
	}()

	select {
	case err := <-errCh:
		fmt.Println(err.Error())
		return 1
	case <-ctx.Done():
		// 終了処理を書く
		fmt.Println("shutting down...")
		return 0
	}
}

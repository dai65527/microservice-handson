package grpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dai65527/microservice-handson/pkg/logger"
	"github.com/dai65527/microservice-handson/platform/db/grpc"
	"github.com/dai65527/microservice-handson/platform/db/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TestDB(t *testing.T) {
	port := 5000

	// サーバ起動
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startServer(t, ctx, port)

	// クライアントの作成
	cli := newClient(t, port)

	// test create customer
	createCustomerRequest := proto.CreateCustomerRequest{
		Name: "Bunjiro",
	}
	createCustomerResponse, err := cli.CreateCustomer(ctx, &createCustomerRequest)
	st, ok := status.FromError(err)
	require.Nil(t, st)
	require.True(t, ok)
	require.NoError(t, err)
	assert.NotEmpty(t, createCustomerResponse.Customer)
	assert.NotEmpty(t, createCustomerResponse.Customer.Id)
	assert.Equal(t, "Bunjiro", createCustomerResponse.Customer.Name)

	bunjiro := *createCustomerResponse.Customer

	createItemRequest := proto.CreateItemRequest{
		CustomerId: bunjiro.Id,
		Title:      "bunjiro marugari",
		Price:      1000,
	}

	createItemResponse, err := cli.CreateItem(ctx, &createItemRequest)
	require.Nil(t, st)
	require.True(t, ok)
	require.NoError(t, err)
	assert.NotNil(t, createItemResponse.Item)
	assert.NotNil(t, createItemResponse.Item.Id)
	assert.Equal(t, bunjiro.Id, createItemResponse.Item.CustomerId)
	assert.Equal(t, "bunjiro marugari", createItemResponse.Item.Title)
	assert.Equal(t, int64(1000), createItemResponse.Item.Price)
}

func startServer(t *testing.T, ctx context.Context, port int) {
	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		require.NoError(t, ferr)
		return
	}
	clogger := l.WithName("db")

	go grpc.RunServer(ctx, 5000, clogger.WithName("grpc"))

	// select {
	// case err := <-errCh:
	// 	fmt.Println(err.Error())
	// 	return err
	// case <-ctx.Done():
	// 	// 終了処理を書く
	// 	fmt.Println("shutting down...")
	// 	return err
	// }
}

func newClient(t *testing.T, port int) proto.DBServiceClient {
	target := fmt.Sprintf("localhost: %d", port)
	conn, err := ggrpc.Dial(target, ggrpc.WithInsecure())
	require.NoError(t, err)
	return proto.NewDBServiceClient(conn)
}

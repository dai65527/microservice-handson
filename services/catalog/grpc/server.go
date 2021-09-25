package grpc

import (
	"context"

	"github.com/dai65527/microservice-handson/services/catalog/proto"
	customer "github.com/dai65527/microservice-handson/services/customer/proto"
	item "github.com/dai65527/microservice-handson/services/item/proto"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.CatalogServiceServer = (*server)(nil)

type server struct {
	proto.UnimplementedCatalogServiceServer
	itemClient     item.ItemServiceClient
	customerClient customer.CustomerServiceClient
}

func (s *server) CreateItem(ctx context.Context, req *proto.CreateItemRequest) (*proto.CreateItemResponse, error) {
	tokenStr, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse access token")
	}

	token, err := jwt.Parse([]byte(tokenStr))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse access token")
	}

	res, err := s.itemClient.CreateItem(ctx, &item.CreateItemRequest{
		CustomerId: token.Subject(),
		Title:      req.Title,
		Price:      req.Price,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	item := res.GetItem()

	return &proto.CreateItemResponse{
		Item: &proto.Item{
			Id:         item.Id,
			CustomerId: item.CustomerId,
			Title:      item.Title,
			Price:      item.Price,
		},
	}, nil
}

func (s *server) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {
	ires, err := s.itemClient.GetItem(ctx, &item.GetItemRequest{Id: req.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	i := ires.GetItem()
	if i == nil {
		return nil, status.Error(codes.NotFound, "internal error")
	}

	cres, err := s.customerClient.GetCustomer(ctx, &customer.GetCustomerRequest{Id: i.CustomerId})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	c := cres.GetCustomer()
	if c == nil {
		return nil, status.Error(codes.NotFound, "internal error")
	}

	return &proto.GetItemResponse{
		Item: &proto.Item{
			Id:           i.Id,
			CustomerId:   i.CustomerId,
			CustomerName: c.Name,
			Title:        i.Title,
			Price:        i.Price,
		},
	}, nil
}

func (s *server) ListItems(ctx context.Context, req *proto.ListItemsRequest) (*proto.ListItemsResponse, error) {
	result, err := s.itemClient.ListItems(ctx, &item.ListItemsRequest{})
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &proto.ListItemsResponse{
		Items: make([]*proto.Item, len(result.GetItems())),
	}

	for i, item := range result.Items {
		res.Items[i] = &proto.Item{
			Id:    item.Id,
			Title: item.Title,
			Price: item.Price,
		}
	}

	return res, nil
}

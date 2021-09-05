package grpc

import (
	"context"

	db "github.com/dai65527/microservice-handson/platform/db/proto"
	"github.com/dai65527/microservice-handson/services/customer/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.CustomerServiceServer = (*server)(nil)

type server struct {
	proto.UnimplementedCustomerServiceServer
	dbClient db.DBServiceClient
}

func (s *server) CreateCustomer(ctx context.Context, req *proto.CreateCustomerRequest) (*proto.CreateCustomerResponse, error) {
	res, err := s.dbClient.CreateCustomer(ctx, &db.CreateCustomerRequest{Name: req.Name})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	customer := res.GetCustomer()

	return &proto.CreateCustomerResponse{
		Customer: &proto.Customer{
			Id:   customer.Id,
			Name: customer.Name,
		},
	}, nil
}

func (s *server) GetCustomer(ctx context.Context, req *proto.GetCustomerRequest) (*proto.GetCustomerResponse, error) {
	res, err := s.dbClient.GetCustomer(ctx, &db.GetCustomerRequest{Id: req.Id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	customer := res.GetCustomer()

	return &proto.GetCustomerResponse{
		Customer: &proto.Customer{
			Id:   customer.Id,
			Name: customer.Name,
		},
	}, nil
}

func (s *server) GetCustomerByName(ctx context.Context, req *proto.GetCustomerByNameRequest) (*proto.GetCustomerByNameResponse, error) {
	res, err := s.dbClient.GetCustomerByName(ctx, &db.GetCustomerByNameRequest{Name: req.Name})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.AlreadyExists, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	customer := res.GetCustomer()

	return &proto.GetCustomerByNameResponse{
		Customer: &proto.Customer{
			Id:   customer.Id,
			Name: customer.Name,
		},
	}, nil
}

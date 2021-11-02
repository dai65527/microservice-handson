package grpc

import (
	"bytes"
	"context"
	"fmt"

	grpccontext "github.com/dai65527/microservice-handson/pkg/grpc/context"
	authority "github.com/dai65527/microservice-handson/services/authority/proto"
	catalog "github.com/dai65527/microservice-handson/services/catalog/proto"
	"github.com/dai65527/microservice-handson/services/gateway/proto"
	"github.com/go-logr/logr"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ proto.GatewayServiceServer   = (*server)(nil)
	_ auth.ServiceAuthFuncOverride = (*server)(nil)

	publicRPCMethods = map[string]struct{}{
		"/dnakano.microservice_handson.gateway.GatewayService/Signup": {},
		"/dnakano.microservice_handson.gateway.GatewayService/Signin": {},
	}
)

type server struct {
	proto.UnimplementedGatewayServiceServer

	authorityClient authority.AuthorityServiceClient
	catalogClient   catalog.CatalogServiceClient
	logger          logr.Logger
}

func (s *server) Signup(ctx context.Context, req *authority.SignupRequest) (*authority.SignupResponse, error) {
	return s.authorityClient.Signup(ctx, req)
}

func (s *server) Signin(ctx context.Context, req *authority.SigninRequest) (*authority.SigninResponse, error) {
	return s.authorityClient.Signin(ctx, req)
}

func (s *server) CreateItem(ctx context.Context, req *catalog.CreateItemRequest) (*catalog.CreateItemResponse, error) {
	return s.catalogClient.CreateItem(ctx, req)
}

func (s *server) GetItem(ctx context.Context, req *catalog.GetItemRequest) (*catalog.GetItemResponse, error) {
	return s.catalogClient.GetItem(ctx, req)
}

func (s *server) ListItems(ctx context.Context, req *catalog.ListItemsRequest) (*catalog.ListItemsResponse, error) {
	return s.catalogClient.ListItems(ctx, req)
}

func (s *server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	_, ok := publicRPCMethods[fullMethodName]
	if ok {
		return ctx, nil
	}

	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		s.log(ctx).Info("failed to get token from authorization header")
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	fmt.Printf("token: %v\n", token)

	res, err := s.authorityClient.ListPublicKeys(ctx, &authority.ListPublicKeysRequest{})
	if err != nil {
		s.log(ctx).Error(err, "failed to call authority's ListPublicKeys")
		return nil, status.Error(codes.Internal, "failed to authenticate")
	}
	fmt.Printf("res.Jwks: %v\n", res.Jwks)

	// key, err := jwk.Parse([]byte(res.Jwks))
	key, err := jwk.Parse(bytes.NewBufferString(res.Jwks).Bytes())
	if err != nil {
		s.log(ctx).Error(err, "failed to parse jwks")
		return nil, status.Error(codes.Internal, "failed to authenticate")
	}
	fmt.Printf("key.Len(): %v\n", key.Len())

	_, err = jwt.Parse([]byte(token), jwt.WithKeySet(key))
	if err != nil {
		s.log(ctx).Info(fmt.Sprintf("failed to verify token: %s", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	return ctx, nil
}

func (s *server) log(ctx context.Context) logr.Logger {
	reqid := grpccontext.GetRequestID(ctx)

	return s.logger.WithValues("request_id", reqid)
}

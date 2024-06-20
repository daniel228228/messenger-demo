package grpc_handler

import (
	"context"

	api "messenger.api/go/api"
)

func (g *grpcHandler) Init(ctx context.Context, req *api.InitRequest) (*api.InitResponse, error) {
	return g.service.Auth().Init(ctx, req)
}

func (g *grpcHandler) Refresh(ctx context.Context, req *api.RefreshRequest) (*api.RefreshResponse, error) {
	return g.service.Auth().Refresh(ctx, req)
}

func (g *grpcHandler) Verify(ctx context.Context, req *api.VerifyRequest) (*api.VerifyResponse, error) {
	return g.service.Auth().Verify(ctx, req)
}

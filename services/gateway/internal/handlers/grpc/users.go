package grpc_handler

import (
	"context"

	"messenger.api/go/api"
)

func (g *grpcHandler) WhoAmI(ctx context.Context, req *api.WhoAmIRequest) (*api.WhoAmIResponse, error) {
	return g.service.Users().WhoAmI(ctx, req)
}

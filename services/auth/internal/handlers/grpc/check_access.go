package grpc_handler

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"messenger.api/go/api"

	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/service"
)

func (g *grpcHandler) CheckAccess(_ context.Context, req *api.CheckAccessRequest) (*api.CheckAccessResponse, error) {
	result, err := g.service.CheckAccess(&dto.AccessToken{
		AccessToken: req.AccessToken,
	})

	switch {
	case errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrExpiredToken):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrInternalError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.CheckAccess").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.CheckAccessResponse{
			UserId: result,
		}, nil
	}
}

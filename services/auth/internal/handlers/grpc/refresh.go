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

func (g *grpcHandler) Refresh(_ context.Context, req *api.RefreshRequest) (*api.RefreshResponse, error) {
	result, err := g.service.RefreshTokens(&dto.RefreshToken{
		RefreshToken: req.RefreshToken,
	})

	switch {
	case errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrExpiredToken):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrInternalError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.Refresh").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.RefreshResponse{
			Token: &api.Token{
				AccessToken:  result.AccessToken,
				RefreshToken: result.RefreshToken,
				ExpiresIn:    int32(result.ExpiresIn),
			},
		}, nil
	}
}

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

func (g *grpcHandler) Verify(_ context.Context, req *api.VerifyRequest) (*api.VerifyResponse, error) {
	verify := &dto.Verify{
		Phone:    req.Phone,
		Password: req.Code,
	}

	result, err := g.service.Verify(verify)

	switch {
	case errors.Is(err, service.ErrInvalidPhoneNumber), errors.Is(err, service.ErrIncorrectCode):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrInternalError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.Verify").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.VerifyResponse{
			Token: &api.Token{
				AccessToken:  result.AccessToken,
				RefreshToken: result.RefreshToken,
				ExpiresIn:    int32(result.ExpiresIn),
			},
		}, nil
	}
}

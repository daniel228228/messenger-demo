package grpc_handler

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"messenger.api/go/api"

	"messenger.users/internal/models/dto"
	"messenger.users/internal/service"
)

func (g *grpcHandler) GetUser(_ context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	result, err := g.service.GetUser(&dto.GetUser{
		ID: req.UserId,
	})

	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrInvalidUserID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrInternalError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.GetUser").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.GetUserResponse{
			User: &api.User{
				Id:        result.ID,
				Username:  result.Username,
				FirstName: result.FirstName,
				LastName:  result.LastName,
			},
		}, nil
	}
}

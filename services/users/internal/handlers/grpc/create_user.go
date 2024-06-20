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

func (g *grpcHandler) CreateUser(_ context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	result, err := g.service.CreateUser(&dto.CreateUser{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	switch {
	case errors.Is(err, service.ErrInternalError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.GetUser").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.CreateUserResponse{
			UserId: result.ID,
		}, nil
	}
}

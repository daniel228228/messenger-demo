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

func (g *grpcHandler) Init(_ context.Context, req *api.InitRequest) (*api.InitResponse, error) {
	err := g.service.InitVerify(&dto.InitVerify{
		Phone: req.Phone,
	})

	switch {
	case errors.Is(err, service.ErrInvalidPhoneNumber):
		return nil, status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrOtpServiceError):
		return nil, status.Error(codes.Unavailable, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.Init").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return &api.InitResponse{}, nil
	}
}

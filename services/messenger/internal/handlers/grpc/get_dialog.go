package grpc_handler

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"messenger.api/go/api"

	"messenger.messenger/internal/service"
)

func (g *grpcHandler) GetDialog(ctx context.Context, req *api.GetDialogRequest) (*api.GetDialogResponse, error) {
	userID, err := g.ExtractUserID(ctx)
	if err != nil {
		g.log.Error().Str("from", "grpcHandler.GetDialog").Msg("Extract UserID error")
		return nil, err
	}

	result, err := g.service.GetDialog(userID, req)

	switch {
	case errors.Is(err, service.ErrDialogNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.GetDialog").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return result, nil
	}
}

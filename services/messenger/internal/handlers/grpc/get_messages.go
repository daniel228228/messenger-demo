package grpc_handler

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"messenger.api/go/api"

	"messenger.messenger/internal/service"
)

func (g *grpcHandler) GetMessages(ctx context.Context, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error) {
	userID, err := g.ExtractUserID(ctx)
	if err != nil {
		g.log.Error().Str("from", "grpcHandler.GetMessages").Msg("Extract UserID error")
		return nil, err
	}

	result, err := g.service.GetMessages(userID, req)

	switch {
	case errors.Is(err, service.ErrBadID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.GetMessages").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return result, nil
	}
}

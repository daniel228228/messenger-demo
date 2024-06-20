package grpc_handler

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"messenger.api/go/api"

	"messenger.messenger/internal/service"
)

func (g *grpcHandler) GetUnreadDialogsCounter(ctx context.Context, req *api.GetUnreadDialogsCounterRequest) (*api.GetUnreadDialogsCounterResponse, error) {
	userID, err := g.ExtractUserID(ctx)
	if err != nil {
		g.log.Error().Str("from", "grpcHandler.GetUnreadDialogsCounter").Msg("Extract UserID error")
		return nil, err
	}

	result, err := g.service.GetUnreadDialogsCounter(userID, req)

	switch {
	case errors.Is(err, service.ErrBadID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		g.log.Error().Str("from", "grpcHandler.GetUnreadDialogsCounter").Err(err).Send()
		return nil, status.Error(codes.Unavailable, "unknown")
	default:
		return result, nil
	}
}

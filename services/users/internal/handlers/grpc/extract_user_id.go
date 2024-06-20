package grpc_handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (g *grpcHandler) ExtractUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		g.log.Error().Str("from", "grpcHandler.ExtractUserID").Msg("missing context metadata")
		return "", status.Error(codes.Internal, "unknown")
	}

	userIDHeader, ok := md["user_id"]
	if !ok {
		g.log.Error().Str("from", "grpcHandler.ExtractUserID").Msg("missing user_id")
		return "", status.Error(codes.Internal, "unknown")
	}

	return userIDHeader[0], nil
}

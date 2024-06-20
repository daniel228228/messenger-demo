package grpc_handler

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *grpcHandler) interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	var userID string

	switch info.FullMethod {
	case // list of methods where authorization is not required
		"/api.auth.Auth/Init",
		"/api.auth.Auth/Verify",
		"/api.auth.Auth/Refresh":
	default:
		var err error

		if ctx, userID, err = g.authInterceptor(ctx); err != nil { // replace context by outgoing ctx with metadata (user_id)
			return nil, err
		}
	}

	hand, err := handler(ctx, req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unavailable: // code conversion
				return nil, status.Error(codes.Unavailable, "service is temporarily unavailable")
			default:
			}
		} else {
			g.log.Error().Err(err).Str("from", "grpcHandler.interceptor").Msg("error parsing handler func status")
			return nil, status.Error(codes.Internal, "unknown")
		}
	}

	g.log.Info().Str("method", info.FullMethod).Str("user_id", userID).Dur("duration", time.Since(start)).Err(err).Send()

	return hand, err
}

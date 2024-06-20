package grpc_handler

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"messenger.api/go/api"
)

func (g *grpcHandler) authInterceptor(ctx context.Context) (context.Context, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, "", status.Error(codes.InvalidArgument, "retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, "", status.Error(codes.Unauthenticated, "access token is not supplied")
	}

	if authHeader[0] == "" {
		return nil, "", status.Error(codes.Unauthenticated, "access token is empty")
	}

	access := strings.Split(authHeader[0], " ")
	if len(access) != 2 || access[0] != "Bearer" {
		return nil, "", status.Error(codes.Unauthenticated, "unsupported type of access token")
	}

	token := access[1]

	result, err := g.service.Auth().CheckAccess(context.Background(), &api.CheckAccessRequest{
		AccessToken: token,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.PermissionDenied: // code conversion
				return nil, "", status.Error(codes.Unauthenticated, e.Message())
			case codes.Unavailable:
				g.log.Error().Err(err).Str("from", "gprcHandler.authInterceptor").Msg("Unavailable Auth.CheckAccess service")
				return nil, "", status.Error(codes.Unavailable, "service is temporarily unavailable")
			default:
				return nil, "", err
			}
		} else {
			g.log.Error().Err(err).Str("from", "gprcHandler.authInterceptor").Msg("Errors parsing Auth.CheckAccess status")
			return nil, "", status.Error(codes.Internal, "unknown")
		}
	}

	md.Append("user_id", result.UserId)

	// creating outgoint context because it's need to be passed
	// user_id metadata field to other services
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx, result.UserId, nil
}

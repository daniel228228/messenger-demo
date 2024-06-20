package grpc_handler

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"messenger.api/go/api"
	"messenger.gateway/pkg/log"

	"messenger.gateway/internal/service"
	"messenger.gateway/pkg/config"
)

var Name = "GrpcHandler"

type GrpcHandler interface{}

type grpcHandler struct {
	config  *config.Config
	log     *log.Logger
	service service.Service
	server  *grpc.Server

	api.UnimplementedAuthServer
	api.UnimplementedUsersServer
	api.UnimplementedMessengerServer
}

func NewGrpcHandler(config *config.Config, log *log.Logger, service service.Service) *grpcHandler {
	return &grpcHandler{
		config:  config,
		log:     log,
		service: service,
	}
}

func (g *grpcHandler) Start(_ context.Context) error {
	g.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	lst, err := net.Listen("tcp", ":"+g.config.GRPCPort)
	if err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to listen")
		return err
	}

	g.server = grpc.NewServer(
		grpc.UnaryInterceptor(g.interceptor),
	)

	api.RegisterAuthServer(g.server, g)
	api.RegisterUsersServer(g.server, g)
	api.RegisterMessengerServer(g.server, g)

	if err := g.server.Serve(lst); err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to serve")
		return err
	}

	return nil
}

func (g *grpcHandler) Stop(_ context.Context) error {
	g.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	g.server.GracefulStop()

	return nil
}

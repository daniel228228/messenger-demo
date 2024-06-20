package grpc_handler

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"messenger.api/go/api"
	"messenger.users/pkg/log"

	"messenger.users/internal/service"
	"messenger.users/pkg/config"
)

var Name = "GrpcHandler"

type GrpcHandler interface{}

type unimplementedUsersServer = api.UnimplementedUsersServer

type grpcHandler struct {
	unimplementedUsersServer
	config  *config.Config
	log     *log.Logger
	service service.Service
	server  *grpc.Server
}

func NewGrpcHandler(config *config.Config, log *log.Logger, service service.Service) *grpcHandler {
	return &grpcHandler{
		config:  config,
		log:     log,
		service: service,
	}
}

func (g *grpcHandler) Start(ctx context.Context) error {
	g.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	lst, err := net.Listen("tcp", ":"+g.config.Port)
	if err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to listen")
		return err
	}

	g.server = grpc.NewServer([]grpc.ServerOption{}...)
	api.RegisterUsersServer(g.server, g)

	if err := g.server.Serve(lst); err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to serve")
	}

	return nil
}

func (g *grpcHandler) Stop(_ context.Context) error {
	g.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	g.server.GracefulStop()

	return nil
}

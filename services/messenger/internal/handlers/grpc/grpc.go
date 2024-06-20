package grpc_handler

import (
	"context"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"messenger.api/go/api"

	"messenger.messenger/pkg/log"

	"messenger.messenger/pkg/config"

	"messenger.messenger/internal/service"
)

var Name = "GrpcHandler"

type GrpcHandler interface{}

type grpcHandler struct {
	api.UnimplementedMessengerServer
	config   *config.Config
	log      *log.Logger
	service  service.Service
	server   *grpc.Server
	shutdown int32
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

	lst, err := net.Listen("tcp", ":"+g.config.Port)
	if err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to listen")
		return err
	}

	g.server = grpc.NewServer()

	api.RegisterMessengerServer(g.server, g)

	if err := g.server.Serve(lst); err != nil {
		g.log.Error().Err(err).Str("from", "grpcHandler").Msg("failed to serve")
		return err
	}

	return nil
}

func (g *grpcHandler) Stop(_ context.Context) error {
	g.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	atomic.StoreInt32(&g.shutdown, 1)

	g.server.GracefulStop()

	return nil
}

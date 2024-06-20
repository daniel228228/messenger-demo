package http_handler

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"messenger.api/go/api"

	"messenger.gateway/pkg/config"
	"messenger.gateway/pkg/log"
)

var Name = "GrpcHandler"

type HTTPHandler interface{}

type httpHandler struct {
	config *config.Config
	log    *log.Logger
	server *http.Server
}

func NewHTTPHandler(config *config.Config, log *log.Logger) *httpHandler {
	return &httpHandler{
		config: config,
		log:    log,
	}
}

func (h *httpHandler) Start(_ context.Context) error {
	h.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	conn, err := grpc.Dial(":"+h.config.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		h.log.Error().Err(err).Str("from", "httpHandler").Msg("can't connect to gRPC server")
		return err
	}
	defer conn.Close()

	mux := runtime.NewServeMux()

	if err = api.RegisterAuthHandler(context.Background(), mux, conn); err != nil {
		h.log.Error().Err(err).Str("from", "httpHandler").Msg("can't register gRPC server")
		return err
	}
	if err = api.RegisterUsersHandler(context.Background(), mux, conn); err != nil {
		h.log.Error().Err(err).Str("from", "httpHandler").Msg("can't register gRPC server")
		return err
	}
	if err = api.RegisterMessengerHandler(context.Background(), mux, conn); err != nil {
		h.log.Error().Err(err).Str("from", "httpHandler").Msg("can't register gRPC server")
		return err
	}

	h.server = &http.Server{
		Addr:    ":" + h.config.Port,
		Handler: mux,
	}

	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.log.Error().Err(err).Str("from", "httpHandler").Msg("ListenAndServe() error")
		return err
	}

	return nil
}

func (h *httpHandler) Stop(ctx context.Context) error {
	h.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return h.server.Shutdown(ctx)
}

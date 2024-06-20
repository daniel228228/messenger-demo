package metrics_handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"messenger.users/pkg/config"
	"messenger.users/pkg/log"
)

var Name = "MetricsHandler"

type MetricsHandler interface{}

type metricsHandler struct {
	config *config.Config
	log    *log.Logger
	server *http.Server
}

func NewMetricsHandler(config *config.Config, log *log.Logger) *metricsHandler {
	return &metricsHandler{
		config: config,
		log:    log,
	}
}

func (m *metricsHandler) Start(_ context.Context) error {
	m.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	m.server = &http.Server{
		Addr:    ":" + m.config.MetricsPort,
		Handler: router,
	}

	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		m.log.Error().Err(err).Str("from", "metricsHandler").Msg("ListenAndServe() error")
		return err
	}

	return nil
}

func (m *metricsHandler) Stop(ctx context.Context) error {
	m.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return m.server.Shutdown(ctx)
}

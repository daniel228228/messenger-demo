package app

import (
	"context"
	"runtime/debug"

	"github.com/pkg/errors"

	grpc_handler "messenger.users/internal/handlers/grpc"
	metrics_handler "messenger.users/internal/handlers/metrics"
	"messenger.users/internal/repo"
	"messenger.users/internal/service"
	"messenger.users/pkg/config"

	"messenger.users/pkg/log"
)

var Name = "users"

type Control interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type component[T any] struct {
	Control
	Name string
	Impl T
}

type application struct {
	Repo           component[repo.Repo]
	Service        component[service.Service]
	GrpcHandler    component[grpc_handler.GrpcHandler]
	MetricsHandler component[metrics_handler.MetricsHandler]
}

var App *application

func init() {
	App = &application{}
}

func Run(ctx context.Context, config *config.Config, log *log.Logger) {
	// repo
	select {
	case <-ctx.Done():
		return
	default:
	}
	repoImpl := repo.NewRepo(config, log)
	App.Repo.Name = repo.Name
	App.Repo.Control = repoImpl
	App.Repo.Impl = repoImpl

	if !start(ctx, App.Repo) {
		return
	}
	defer stop(ctx, App.Repo)

	// service
	select {
	case <-ctx.Done():
		return
	default:
	}
	serviceImpl := service.NewService(config, log, App.Repo.Impl)
	App.Service.Name = service.Name
	App.Service.Control = serviceImpl
	App.Service.Impl = serviceImpl

	if !start(ctx, App.Service) {
		return
	}
	defer stop(ctx, App.Service)

	// handler
	select {
	case <-ctx.Done():
		return
	default:
	}
	grpcHandlerImpl := grpc_handler.NewGrpcHandler(config, log, App.Service.Impl)
	App.GrpcHandler.Name = grpc_handler.Name
	App.GrpcHandler.Control = grpcHandlerImpl
	App.GrpcHandler.Impl = grpcHandlerImpl

	start_parallel(ctx, log, App.GrpcHandler)
	defer stop(ctx, App.GrpcHandler)

	// metrics_handler
	select {
	case <-ctx.Done():
		return
	default:
	}
	metricsHandlerImpl := metrics_handler.NewMetricsHandler(config, log)
	App.MetricsHandler.Name = metrics_handler.Name
	App.MetricsHandler.Control = metricsHandlerImpl
	App.MetricsHandler.Impl = metricsHandlerImpl

	start_parallel(ctx, log, App.MetricsHandler)
	defer stop(ctx, App.MetricsHandler)

	<-ctx.Done()
}

func start_parallel[T any](ctx context.Context, log *log.Logger, c component[T]) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal().Bool("app", true).Msgf("panic: %v\n\n%s", r, string(debug.Stack()))
			}
		}()

		_ = c.Start(ctx)
	}()
}

func start[T any](ctx context.Context, c component[T]) bool {
	if err := c.Start(ctx); err != nil {
		return false
	}

	return true
}

func stop[T any](ctx context.Context, c component[T]) {
	if err := c.Stop(ctx); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		// nop
	}
}

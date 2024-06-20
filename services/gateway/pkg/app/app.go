package app

import (
	"context"
	"runtime/debug"

	"github.com/pkg/errors"

	"messenger.gateway/internal/grpc_client"
	grpc_handler "messenger.gateway/internal/handlers/grpc"
	http_handler "messenger.gateway/internal/handlers/http"
	metrics_handler "messenger.gateway/internal/handlers/metrics"
	"messenger.gateway/internal/service"
	"messenger.gateway/pkg/config"

	"messenger.gateway/pkg/log"
)

var Name = "gateway"

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
	Auth      component[grpc_client.Auth]
	Users     component[grpc_client.Users]
	Messenger component[grpc_client.Messenger]

	Service        component[service.Service]
	GrpcHandler    component[grpc_handler.GrpcHandler]
	HTTPHandler    component[http_handler.HTTPHandler]
	MetricsHandler component[metrics_handler.MetricsHandler]
}

var App *application

func init() {
	App = &application{}
}

func Run(ctx context.Context, config *config.Config, log *log.Logger) {
	// auth
	select {
	case <-ctx.Done():
		return
	default:
	}
	authImpl := grpc_client.NewAuth(config, log)
	App.Auth.Name = authImpl.Name()
	App.Auth.Control = authImpl
	App.Auth.Impl = authImpl

	start_parallel(ctx, log, App.Auth)
	defer stop(ctx, App.Auth)

	// users
	select {
	case <-ctx.Done():
		return
	default:
	}
	usersImpl := grpc_client.NewUsers(config, log)
	App.Users.Name = usersImpl.Name()
	App.Users.Control = usersImpl
	App.Users.Impl = usersImpl

	start_parallel(ctx, log, App.Users)
	defer stop(ctx, App.Users)

	// messenger
	select {
	case <-ctx.Done():
		return
	default:
	}
	messengerImpl := grpc_client.NewMessenger(config, log)
	App.Messenger.Name = messengerImpl.Name()
	App.Messenger.Control = messengerImpl
	App.Messenger.Impl = messengerImpl

	start_parallel(ctx, log, App.Messenger)
	defer stop(ctx, App.Messenger)

	// service
	select {
	case <-ctx.Done():
		return
	default:
	}
	serviceImpl := service.NewService(config, log, App.Auth.Impl, App.Users.Impl, App.Messenger.Impl)
	App.Service.Name = service.Name
	App.Service.Control = serviceImpl
	App.Service.Impl = serviceImpl

	if !start(ctx, App.Service) {
		return
	}
	defer stop(ctx, App.Service)

	// grpc_handler
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

	// http_handler
	select {
	case <-ctx.Done():
		return
	default:
	}
	httpHandlerImpl := http_handler.NewHTTPHandler(config, log)
	App.HTTPHandler.Name = http_handler.Name
	App.HTTPHandler.Control = httpHandlerImpl
	App.HTTPHandler.Impl = httpHandlerImpl

	start_parallel(ctx, log, App.HTTPHandler)
	defer stop(ctx, App.HTTPHandler)

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

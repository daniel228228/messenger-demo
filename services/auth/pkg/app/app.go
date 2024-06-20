package app

import (
	"context"
	"runtime/debug"

	"github.com/pkg/errors"

	"messenger.auth/internal/cache"
	"messenger.auth/internal/grpc_client"
	grpc_handler "messenger.auth/internal/handlers/grpc"
	metrics_handler "messenger.auth/internal/handlers/metrics"
	"messenger.auth/internal/otp_service"
	"messenger.auth/internal/service"
	"messenger.auth/pkg/config"

	"messenger.auth/pkg/log"
)

var Name = "auth"

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
	Cache          component[cache.Cache]
	OtpService     component[otp_service.OtpService]
	Users          component[grpc_client.Users]
	Service        component[service.Service]
	GrpcHandler    component[grpc_handler.GrpcHandler]
	MetricsHandler component[metrics_handler.MetricsHandler]
}

var App *application

func init() {
	App = &application{}
}

func Run(ctx context.Context, config *config.Config, log *log.Logger) {
	// cache
	select {
	case <-ctx.Done():
		return
	default:
	}
	cacheImpl := cache.NewCache(config, log)
	App.Cache.Name = cache.Name
	App.Cache.Control = cacheImpl
	App.Cache.Impl = cacheImpl

	if !start(ctx, App.Cache) {
		return
	}
	defer stop(ctx, App.Cache)

	// otp_service
	select {
	case <-ctx.Done():
		return
	default:
	}

	otpServiceImpl := otp_service.NewMockOtpService(config, log)
	App.OtpService.Name = otp_service.Name
	App.OtpService.Control = otpServiceImpl
	App.OtpService.Impl = otpServiceImpl

	if !start(ctx, App.OtpService) {
		return
	}
	defer stop(ctx, App.OtpService)

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

	// service
	select {
	case <-ctx.Done():
		return
	default:
	}
	serviceImpl := service.NewService(config, log, App.Cache.Impl, App.OtpService.Impl, App.Users.Impl)
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

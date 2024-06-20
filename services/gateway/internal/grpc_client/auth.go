package grpc_client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"messenger.api/go/api"

	"messenger.gateway/pkg/log"

	"messenger.gateway/pkg/config"
)

type Auth interface {
	api.AuthClient
}

type auth struct {
	cl     api.AuthClient
	conn   *grpc.ClientConn
	config *config.Config
	log    *log.Logger
}

func NewAuth(config *config.Config, log *log.Logger) *auth {
	return &auth{
		config: config,
		log:    log,
	}
}

func (cl *auth) Start(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "start").Send()

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	cl.conn, err = grpc.DialContext(
		ctx,
		cl.config.AuthUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Duration(1 * time.Second),
				Multiplier: backoff.DefaultConfig.Multiplier,
				Jitter:     backoff.DefaultConfig.Jitter,
				MaxDelay:   time.Duration(120 * time.Second),
			},
			MinConnectTimeout: time.Duration(1 * time.Second),
		}),
	)
	if err != nil {
		return err
	}

	cl.cl = api.NewAuthClient(cl.conn)

	return nil
}

func (cl *auth) Stop(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "stop").Send()

	if cl.conn != nil {
		cl.conn.Close()
	}

	return nil
}

func (cl *auth) Name() string {
	return "grpcClient (Auth)"
}

func (cl *auth) CheckAccess(ctx context.Context, in *api.CheckAccessRequest, opts ...grpc.CallOption) (*api.CheckAccessResponse, error) {
	return cl.cl.CheckAccess(ctx, in, opts...)
}

func (cl *auth) Init(ctx context.Context, in *api.InitRequest, opts ...grpc.CallOption) (*api.InitResponse, error) {
	return cl.cl.Init(ctx, in, opts...)
}

func (cl *auth) Verify(ctx context.Context, in *api.VerifyRequest, opts ...grpc.CallOption) (*api.VerifyResponse, error) {
	return cl.cl.Verify(ctx, in, opts...)
}

func (cl *auth) Refresh(ctx context.Context, in *api.RefreshRequest, opts ...grpc.CallOption) (*api.RefreshResponse, error) {
	return cl.cl.Refresh(ctx, in, opts...)
}

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

type Users interface {
	api.UsersClient
}

type users struct {
	cl     api.UsersClient
	conn   *grpc.ClientConn
	config *config.Config
	log    *log.Logger
}

func NewUsers(config *config.Config, log *log.Logger) *users {
	return &users{
		config: config,
		log:    log,
	}
}

func (cl *users) Start(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "start").Send()

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	cl.conn, err = grpc.DialContext(
		ctx,
		cl.config.UsersUrl,
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

	cl.cl = api.NewUsersClient(cl.conn)

	return nil
}

func (cl *users) Stop(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "stop").Send()

	if cl.conn != nil {
		cl.conn.Close()
	}

	return nil
}

func (cl *users) Name() string {
	return "grpcClient (Users)"
}

func (cl *users) GetUser(ctx context.Context, in *api.GetUserRequest, opts ...grpc.CallOption) (*api.GetUserResponse, error) {
	return nil, nil
}

func (cl *users) CreateUser(ctx context.Context, in *api.CreateUserRequest, opts ...grpc.CallOption) (*api.CreateUserResponse, error) {
	return nil, nil
}

func (cl *users) WhoAmI(ctx context.Context, in *api.WhoAmIRequest, opts ...grpc.CallOption) (*api.WhoAmIResponse, error) {
	return cl.cl.WhoAmI(ctx, in, opts...)
}

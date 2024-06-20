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

type Messenger interface {
	api.MessengerClient
}

type messenger struct {
	cl     api.MessengerClient
	conn   *grpc.ClientConn
	config *config.Config
	log    *log.Logger
}

func NewMessenger(config *config.Config, log *log.Logger) *messenger {
	return &messenger{
		config: config,
		log:    log,
	}
}

func (cl *messenger) Start(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "start").Send()

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	cl.conn, err = grpc.DialContext(
		ctx,
		cl.config.MessengerUrl,
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

	cl.cl = api.NewMessengerClient(cl.conn)

	return nil
}

func (cl *messenger) Stop(ctx context.Context) error {
	cl.log.Info().Bool("app", true).Str("component", cl.Name()).Str("state", "stop").Send()

	if cl.conn != nil {
		cl.conn.Close()
	}

	return nil
}

func (cl *messenger) Name() string {
	return "grpcClient (Messenger)"
}

func (cl *messenger) SendMessage(ctx context.Context, in *api.SendMessageRequest, opts ...grpc.CallOption) (*api.SendMessageResponse, error) {
	return cl.cl.SendMessage(ctx, in, opts...)
}

func (cl *messenger) ReadMessage(ctx context.Context, in *api.ReadMessageRequest, opts ...grpc.CallOption) (*api.ReadMessageResponse, error) {
	return cl.cl.ReadMessage(ctx, in, opts...)
}

func (cl *messenger) GetDialog(ctx context.Context, in *api.GetDialogRequest, opts ...grpc.CallOption) (*api.GetDialogResponse, error) {
	return cl.cl.GetDialog(ctx, in, opts...)
}

func (cl *messenger) GetDialogs(ctx context.Context, in *api.GetDialogsRequest, opts ...grpc.CallOption) (*api.GetDialogsResponse, error) {
	return cl.cl.GetDialogs(ctx, in, opts...)
}

func (cl *messenger) GetMessages(ctx context.Context, in *api.GetMessagesRequest, opts ...grpc.CallOption) (*api.GetMessagesResponse, error) {
	return cl.cl.GetMessages(ctx, in, opts...)
}

func (cl *messenger) GetUnreadDialogsCounter(ctx context.Context, in *api.GetUnreadDialogsCounterRequest, opts ...grpc.CallOption) (*api.GetUnreadDialogsCounterResponse, error) {
	return cl.cl.GetUnreadDialogsCounter(ctx, in, opts...)
}

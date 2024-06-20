package service

import (
	"context"

	"messenger.gateway/pkg/log"

	"messenger.gateway/internal/grpc_client"
	"messenger.gateway/pkg/config"
)

var Name = "Service"

type Service interface {
	Auth() grpc_client.Auth
	Users() grpc_client.Users
	Messenger() grpc_client.Messenger
}

type service struct {
	config    *config.Config
	log       *log.Logger
	auth      grpc_client.Auth
	users     grpc_client.Users
	messenger grpc_client.Messenger
}

func NewService(config *config.Config, log *log.Logger, grpcClients ...any) *service {
	s := &service{
		config: config,
		log:    log,
	}

	for _, v := range grpcClients {
		switch cl := v.(type) {
		case grpc_client.Auth:
			s.auth = cl
		case grpc_client.Users:
			s.users = cl
		case grpc_client.Messenger:
			s.messenger = cl
		}
	}

	return s
}

func (s *service) Start(_ context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	return nil
}

func (s *service) Stop(_ context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return nil
}

func (s *service) Auth() grpc_client.Auth {
	return s.auth
}

func (s *service) Users() grpc_client.Users {
	return s.users
}

func (s *service) Messenger() grpc_client.Messenger {
	return s.messenger
}

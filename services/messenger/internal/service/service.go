package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"messenger.api/go/api"

	"messenger.messenger/pkg/log"

	"messenger.messenger/pkg/config"

	"messenger.messenger/internal/cache"
	"messenger.messenger/internal/grpc_client"
	"messenger.messenger/internal/repo"
)

var Name = "Service"

type Service interface {
	GetDialog(userID string, req *api.GetDialogRequest) (*api.GetDialogResponse, error)
	GetDialogs(userID string, req *api.GetDialogsRequest) (*api.GetDialogsResponse, error)
	GetMessages(userID string, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error)
	GetUnreadDialogsCounter(userID string, req *api.GetUnreadDialogsCounterRequest) (*api.GetUnreadDialogsCounterResponse, error)

	ReadMessage(userID string, req *api.ReadMessageRequest) (*api.ReadMessageResponse, error)
	SendMessage(userID string, req *api.SendMessageRequest) (*api.SendMessageResponse, error)
}

var (
	ErrBadID          = errors.New("bad id")
	ErrUserNotFound   = errors.New("user not found")
	ErrPeerNotFound   = errors.New("peer not found")
	ErrDialogNotFound = errors.New("dialog not found")
)

type service struct {
	config *config.Config
	log    *log.Logger
	repo   repo.Repo
	cache  cache.Cache

	users grpc_client.Users
}

func NewService(config *config.Config, log *log.Logger, repo repo.Repo, cache cache.Cache, grpcClients ...any) *service {
	s := &service{
		config: config,
		log:    log,
		repo:   repo,
		cache:  cache,
	}

	for _, v := range grpcClients {
		switch cl := v.(type) {
		case grpc_client.Users:
			s.users = cl
		}
	}

	initCachedUsers(s) // TODO: temporary solution

	return s
}

func (s *service) Start(ctx context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	go func() { // TODO: temporary solution
		ticker := time.NewTicker(3 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.log.Info().Msg("Clearing cached users list...")
				cachedUsers.clear()
				s.log.Info().Msg("Cached users list has been cleared")
			}
		}
	}()

	return nil
}

func (s *service) Stop(_ context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return nil
}

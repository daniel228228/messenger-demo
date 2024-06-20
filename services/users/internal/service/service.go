package service

import (
	"context"

	"github.com/pkg/errors"

	"messenger.users/pkg/log"

	"messenger.users/internal/models/dto"
	"messenger.users/internal/repo"
	"messenger.users/pkg/config"
)

var Name = "Service"

type Service interface {
	GetUser(user *dto.GetUser) (*dto.User, error)
	CreateUser(user *dto.CreateUser) (*dto.UserID, error)
}

var (
	ErrInternalError = errors.New("internal error")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user id")
)

type service struct {
	config *config.Config
	log    *log.Logger
	repo   repo.Repo
}

func NewService(config *config.Config, log *log.Logger, repo repo.Repo) *service {
	return &service{
		config: config,
		log:    log,
		repo:   repo,
	}
}

func (s *service) Start(_ context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	return nil
}

func (s *service) Stop(_ context.Context) error {
	s.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return nil
}

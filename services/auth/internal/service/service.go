package service

import (
	"context"

	"github.com/pkg/errors"

	"messenger.auth/pkg/log"

	"messenger.auth/internal/cache"
	"messenger.auth/internal/grpc_client"
	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/otp_service"
	"messenger.auth/internal/service/jwt"
	"messenger.auth/pkg/config"
)

var Name = "Service"

type Service interface {
	InitVerify(initVerify *dto.InitVerify) error
	Verify(verify *dto.Verify) (*dto.Tokens, error)
	RefreshTokens(refresh *dto.RefreshToken) (*dto.Tokens, error)
	CheckAccess(access *dto.AccessToken) (string, error)
}

var (
	ErrOtpServiceError    = errors.New("otp service error")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrIncorrectCode      = errors.New("incorrect code")
	ErrInternalError      = errors.New("internal error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUserID      = errors.New("invalid user id")
)

type service struct {
	config       *config.Config
	log          *log.Logger
	cache        cache.Cache
	otpService   otp_service.OtpService
	usersService grpc_client.Users
	jwt          jwt.JWT
}

func NewService(config *config.Config, log *log.Logger, cache cache.Cache, otpService otp_service.OtpService, usersService grpc_client.Users) *service {
	return &service{
		config:       config,
		log:          log,
		cache:        cache,
		otpService:   otpService,
		usersService: usersService,
		jwt:          jwt.NewJWT(config, log, cache),
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

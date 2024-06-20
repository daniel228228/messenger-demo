package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"messenger.auth/internal/cache"

	"messenger.auth/pkg/log"

	"messenger.auth/internal/models/dto"
	"messenger.auth/pkg/config"
)

type JWT interface {
	GenerateTokens(userID string) (*dto.Tokens, error)
	RefreshTokens(refresh *dto.RefreshToken) (*dto.Tokens, error)
	DeleteToken(token string, t string) (string, error)
	ParseToken(token string) (jwt.MapClaims, error)
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type jwtImpl struct {
	config *config.Config
	log    *log.Logger
	cache  cache.Cache
}

func NewJWT(config *config.Config, log *log.Logger, cache cache.Cache) *jwtImpl {
	return &jwtImpl{
		config: config,
		log:    log,
		cache:  cache,
	}
}

func (j *jwtImpl) GenerateTokens(userID string) (*dto.Tokens, error) {
	return j.generateTokens(userID)
}

func (j *jwtImpl) RefreshTokens(refresh *dto.RefreshToken) (*dto.Tokens, error) {
	return j.refreshTokens(refresh)
}

func (j *jwtImpl) DeleteToken(token string, t string) (string, error) {
	return j.deleteToken(token, t)
}

func (j *jwtImpl) ParseToken(token string) (jwt.MapClaims, error) {
	return j.parseToken(token)
}

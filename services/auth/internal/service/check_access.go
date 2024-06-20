package service

import (
	"github.com/pkg/errors"

	"messenger.auth/internal/cache"
	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/service/jwt"
)

func (s *service) CheckAccess(access *dto.AccessToken) (string, error) {
	claims, err := s.jwt.ParseToken(access.AccessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return "", ErrExpiredToken
		}

		s.log.Error().Str("from", "jwt.deleteToken").Err(errors.Wrapf(err, "parse token error")).Send()
		return "", ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	cachedUserID, err := s.cache.Get(uuid)
	if err != nil {
		if !errors.Is(err, cache.ErrEmpty) {
			s.log.Error().Str("from", "service.CheckAccess (get from cache)").Err(errors.Wrapf(err, "cache: Get error")).Send()
			return "", err
		}

		return "", ErrInvalidToken
	}

	if cachedUserID != userID {
		return "", ErrInvalidToken
	}

	return userID, nil
}

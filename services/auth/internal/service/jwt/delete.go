package jwt

import (
	"github.com/pkg/errors"

	"messenger.auth/internal/cache"
)

func (j *jwtImpl) deleteToken(token string, t string) (string, error) {
	claims, err := j.parseToken(token)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return "", ErrExpiredToken
		}

		j.log.Error().Str("from", "jwt.deleteToken").Err(errors.Wrapf(err, "parse token error")).Send()
		return "", ErrInvalidToken
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != t {
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

	if _, err := j.cache.Get(uuid); err != nil {
		if !errors.Is(err, cache.ErrEmpty) {
			j.log.Error().Str("from", "jwt.deleteToken (get from cache)").Err(errors.Wrapf(err, "cache: Get error")).Send()
			return "", err
		}

		return "", ErrInvalidToken
	}

	if err := j.cache.Del(uuid); err != nil {
		if !errors.Is(err, cache.ErrEmpty) {
			j.log.Error().Str("from", "jwt.deleteToken (del from cache)").Err(errors.Wrapf(err, "cache: Del error")).Send()
			return "", err
		}

		return "", ErrInvalidToken
	}

	return userID, nil
}

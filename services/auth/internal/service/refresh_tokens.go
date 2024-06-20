package service

import (
	"github.com/pkg/errors"

	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/service/jwt"
)

func (s *service) RefreshTokens(refresh *dto.RefreshToken) (*dto.Tokens, error) {
	result, err := s.jwt.RefreshTokens(refresh)

	if errors.Is(err, jwt.ErrInvalidToken) {
		return nil, ErrInvalidToken
	} else if errors.Is(err, jwt.ErrExpiredToken) {
		return nil, ErrExpiredToken
	} else if err != nil {
		s.log.Error().Str("from", "RefreshToken").Err(errors.Wrap(err, "jwt.RefreshTokens error")).Send()
		return nil, ErrInternalError
	}

	return result, nil
}

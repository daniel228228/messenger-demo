package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"messenger.users/internal/models/dto"
	"messenger.users/internal/repo"
)

func (s *service) GetUser(user *dto.GetUser) (*dto.User, error) {
	_, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	u, err := s.repo.User(user.ID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			s.log.Warn().Str("from", "service.GetUser (repo.User)").Str("user_id", user.ID).Msg("User not found")
			return nil, ErrUserNotFound
		}

		s.log.Error().Str("from", "service.GetUser (repo.User)").Err(err).Send()
		return nil, ErrInternalError
	}

	return &dto.User{
		ID:        u.ID.String(),
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

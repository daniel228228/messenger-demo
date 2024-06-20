package service

import (
	"messenger.users/internal/models/dto"
	"messenger.users/internal/models/entities"
)

func (s *service) CreateUser(user *dto.CreateUser) (*dto.UserID, error) {
	entUser := &entities.User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	u, err := s.repo.CreateUser(entUser)
	if err != nil {
		s.log.Error().Str("from", "service.CreateUser (repo.User)").Err(err).Send()
		return nil, ErrInternalError
	}

	return &dto.UserID{
		ID: u.String(),
	}, nil
}

package service

import (
	"context"

	"messenger.api/go/api"
)

func (s *service) GetUnreadDialogsCounter(userID string, req *api.GetUnreadDialogsCounterRequest) (*api.GetUnreadDialogsCounterResponse, error) {
	count, err := s.cache.GetUnreadDialogsCounter(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return &api.GetUnreadDialogsCounterResponse{
		Count: int32(count),
	}, nil
}

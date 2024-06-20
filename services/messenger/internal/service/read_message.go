package service

import (
	"context"

	"github.com/pkg/errors"
	"messenger.api/go/api"

	"messenger.messenger/internal/repo"
)

func (s *service) ReadMessage(userID string, req *api.ReadMessageRequest) (*api.ReadMessageResponse, error) {
	var peerID string

	switch v := req.Peer.Peer.(type) {
	case *api.Peer_User:
		peerID = v.User.UserId
	default:
		return nil, errors.New("unsupported peer type")
	}

	lastID := req.LastId

	if err := s.repo.ReadMessages(userID, peerID, lastID); err != nil {
		switch {
		case errors.Is(err, repo.ErrReadRejected):
			return &api.ReadMessageResponse{}, nil
		case errors.Is(err, repo.ErrBadUserID):
			return nil, errors.Wrap(ErrBadID, "user")
		case errors.Is(err, repo.ErrBadPeerID):
			return nil, errors.Wrap(ErrBadID, "peer")
		case errors.Is(err, repo.ErrBadMessageID):
			return nil, errors.Wrap(ErrBadID, "message")
		default:
			s.log.Error().Str("from", "service.ReadMessage (repo.ReadMessages)").Err(err).Send()
			return nil, err
		}
	}

	isUnread, err := s.repo.IsUnreadDialog(userID, req.Peer.GetUser().UserId)
	if err != nil {
		s.log.Error().Str("from", "service.ReadMessage (repo.IsUnreadDialog)").Err(err).Send()
		return nil, err
	} else if isUnread {
		return &api.ReadMessageResponse{}, nil
	}

	_, err = s.cache.DecrUnreadDialogsCounter(context.Background(), userID)
	if err != nil {
		s.log.Error().Str("from", "service.ReadMessage (cachedDB.DecrUnreadDialogsCounter)").Err(err).Send()
		return nil, err
	}

	return &api.ReadMessageResponse{}, nil
}

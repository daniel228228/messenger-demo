package service

import (
	"github.com/pkg/errors"

	"messenger.messenger/internal/repo"
)

func (s *service) getLastReadMessageID(fromID, toID string) (*string, error) {
	id, err := s.repo.LastReadMessageID(fromID, toID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			s.log.Error().Str("from", "service.getLastReadMsgID").Str("from_id", fromID).Str("to_id", toID).Err(errors.Wrap(ErrBadID, "user")).Send()
		case errors.Is(err, repo.ErrBadPeerID):
			s.log.Error().Str("from", "service.getLastReadMsgID").Str("from_id", fromID).Str("to_id", toID).Err(errors.Wrap(ErrBadID, "peer")).Send()
		default:
			s.log.Error().Str("from", "service.getLastReadMsgID").Str("from_id", fromID).Str("to_id", toID).Err(err).Send()
		}

		return nil, err
	}

	return id, nil
}

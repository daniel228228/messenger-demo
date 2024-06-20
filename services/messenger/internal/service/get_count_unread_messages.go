package service

import (
	"github.com/pkg/errors"

	"messenger.messenger/internal/repo"
)

func (s *service) getCountUnreadMessages(fromID, toID string, lastReadMsgID *string) (int, error) {
	count, err := s.repo.CountUnreadMessages(fromID, toID, lastReadMsgID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			s.log.Error().Str("from", "service.getCountUnreadMessages").Str("from_id", fromID).Str("to_id", toID).Err(errors.Wrap(ErrBadID, "user")).Send()
		case errors.Is(err, repo.ErrBadPeerID):
			s.log.Error().Str("from", "service.getCountUnreadMessages").Str("from_id", fromID).Str("to_id", toID).Err(errors.Wrap(ErrBadID, "peer")).Send()
		case errors.Is(err, repo.ErrBadMessageID):
			s.log.Error().Str("from", "service.getCountUnreadMessages").Str("from_id", fromID).Str("to_id", toID).Str("last_message_id", *lastReadMsgID).Err(errors.Wrap(ErrBadID, "message")).Send()
		default:
			s.log.Error().Str("from", "service.getCountUnreadMessages").Str("from_id", fromID).Str("to_id", toID).Str("last_message_id", *lastReadMsgID).Err(err).Send()
		}

		return 0, err
	}

	return count, nil
}

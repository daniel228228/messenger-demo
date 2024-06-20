package service

import (
	"github.com/pkg/errors"
	"messenger.api/go/api"

	"messenger.messenger/internal/models/entities"

	"messenger.messenger/internal/repo"
)

func (s *service) getLastMessage(userID, peerID string) (*api.Message, string) {
	entMessages := []entities.MessageWithAuthor{}

	_, err := s.repo.Messages(userID, peerID, nil, 1, &entMessages)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			s.log.Error().Str("from", "service.getLastMessage (repo.Messages)").Err(errors.Wrap(ErrBadID, "user")).Send()
			return nil, ""
		case errors.Is(err, repo.ErrBadPeerID):
			s.log.Error().Str("from", "service.getLastMessage (repo.Messages)").Err(errors.Wrap(ErrBadID, "peer")).Send()
			return nil, ""
		default:
			s.log.Error().Str("from", "service.getLastMessage (repo.Messages)").Err(err).Send()
			return nil, ""
		}
	}

	if len(entMessages) > 0 {
		return s.buildPeerMessageProto(&entMessages[0].Message, userID == entMessages[0].Author.String(), &api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: peerID,
				},
			},
		}, &api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: entMessages[0].Author.String(),
				},
			},
		}), entMessages[0].ID.String()
	}

	return nil, ""
}

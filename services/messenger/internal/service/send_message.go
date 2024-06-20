package service

import (
	"context"

	"github.com/pkg/errors"

	"messenger.api/go/api"

	"messenger.messenger/internal/models/entities"
	"messenger.messenger/internal/repo"
)

func (s *service) SendMessage(userID string, req *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	entMsg := &entities.Message{
		Message: req.Message,
	}

	var peerID string

	switch v := req.Peer.Peer.(type) {
	case *api.Peer_User:
		peerID = v.User.UserId
	default:
		return nil, errors.New("unsupported peer type")
	}

	wasUnreadDialogByPeer, err := s.repo.IsUnreadDialog(peerID, userID)
	if err != nil {
		switch { // reverse some errors for SendMessage caller client
		case errors.Is(err, repo.ErrBadUserID):
			return nil, errors.Wrap(ErrBadID, "peer")
		case errors.Is(err, repo.ErrBadPeerID):
			return nil, errors.Wrap(ErrBadID, "user")
		default:
			s.log.Error().Str("from", "service.SendMessage (repo.IsUnreadDialog)").Err(err).Send()
			return nil, err
		}
	}

	if err := s.repo.SaveMessage(userID, peerID, entMsg); err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			return nil, errors.Wrap(ErrBadID, "user")
		case errors.Is(err, repo.ErrBadPeerID):
			return nil, errors.Wrap(ErrBadID, "peer")
		case errors.Is(err, repo.ErrPeerNotFound):
			return nil, ErrPeerNotFound
		default:
			s.log.Error().Str("from", "service.SendMessage (repo.SaveMessage)").Err(err).Send()
			return nil, err
		}
	}

	if !wasUnreadDialogByPeer {
		_, err := s.cache.IncrUnreadDialogsCounter(context.Background(), peerID)
		if err != nil {
			s.log.Error().Str("from", "service.SendMessage (cache.IncrUnreadDialogsCounter)").Err(err).Send()
		}
	}

	msg := s.buildPeerMessageProto(entMsg, true, req.Peer, &api.Peer{
		Peer: &api.Peer_User{
			User: &api.PeerUser{
				UserId: userID,
			},
		},
	})

	return &api.SendMessageResponse{
		Response: &api.SendMessageResponse_Message{
			Message: msg,
		},
	}, nil
}

package service

import (
	"github.com/pkg/errors"

	"messenger.api/go/api"

	"messenger.messenger/internal/models/entities"
	"messenger.messenger/internal/repo"
)

func (s *service) GetMessages(userID string, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error) {
	entMessages := []entities.MessageWithAuthor{}

	var peerID string

	switch v := req.Peer.Peer.(type) {
	case *api.Peer_User:
		peerID = v.User.UserId
	default:
		return nil, errors.New("unsupported peer type")
	}

	total, err := s.repo.Messages(userID, peerID, req.OffsetId, int(req.Limit), &entMessages)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			return nil, errors.Wrap(ErrBadID, "user")
		case errors.Is(err, repo.ErrBadPeerID):
			return nil, errors.Wrap(ErrBadID, "peer")
		case errors.Is(err, repo.ErrBadMessageID):
			return nil, errors.Wrap(ErrBadID, "message")
		default:
			s.log.Error().Str("from", "service.GetMessages (repo.Messages)").Err(err).Send()
			return nil, err
		}
	}

	messages := make([]*api.Message, 0, len(entMessages))

	for _, v := range entMessages {
		messages = append(messages, s.buildPeerMessageProto(&v.Message, userID == v.Author.String(), req.Peer, &api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: v.Author.String(),
				},
			},
		}))
	}

	var users []*api.User

	if peerID != userID {
		users = []*api.User{cachedUsers.get(userID), cachedUsers.get(peerID)}
	} else {
		users = []*api.User{cachedUsers.get(userID)}
	}

	return &api.GetMessagesResponse{
		Messages: &api.GetMessagesResponse_Slice{
			Slice: &api.MessagesSlice{
				Total:    int32(total),
				Messages: messages,
				Users:    users,
			},
		},
	}, nil
}

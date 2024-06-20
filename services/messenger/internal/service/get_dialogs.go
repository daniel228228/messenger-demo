package service

import (
	"github.com/pkg/errors"

	"messenger.messenger/internal/models/entities"

	"messenger.messenger/internal/repo"

	"messenger.api/go/api"
)

func (s *service) GetDialogs(userID string, req *api.GetDialogsRequest) (*api.GetDialogsResponse, error) {
	entMessages := []entities.MessageWithAuthorRecipient{}

	total, err := s.repo.Dialogs(userID, req.OffsetDate.AsTime(), int(req.Limit), &entMessages)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrBadUserID):
			return nil, errors.Wrap(ErrBadID, "user")
		default:
			s.log.Error().Str("from", "service.GetDialogs (repo.Dialogs)").Err(err).Send()
			return nil, err
		}
	}

	dialogs := make([]*api.Dialog, 0, len(entMessages))
	messages := make([]*api.Message, 0, len(entMessages))
	users := []*api.User{}

	for _, v := range entMessages {
		var peerID string

		if userID == v.Author.String() {
			peerID = v.Recipient.String()
		} else {
			peerID = v.Author.String()
		}

		lastReadMsgIDFromUser, _ := s.getLastReadMessageID(userID, peerID)
		lastReadMsgIDFromPeer, _ := s.getLastReadMessageID(peerID, userID)
		countUnreadMsgs, _ := s.getCountUnreadMessages(userID, peerID, lastReadMsgIDFromUser)

		dialogs = append(dialogs, s.buildDialogProto(&api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: peerID,
				},
			},
		}, v.ID.String(), countUnreadMsgs, lastReadMsgIDFromUser, lastReadMsgIDFromPeer))

		messages = append(messages, s.buildPeerMessageProto(&v.Message, userID == v.Author.String(), &api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: peerID,
				},
			},
		}, &api.Peer{
			Peer: &api.Peer_User{
				User: &api.PeerUser{
					UserId: v.Author.String(),
				},
			},
		}))

		if !searchID(peerID, users) {
			users = append(users, cachedUsers.get(peerID))
		}
	}

	return &api.GetDialogsResponse{
		Dialogs: &api.GetDialogsResponse_Slice{
			Slice: &api.DialogsSlice{
				Total:    int32(total),
				Dialogs:  dialogs,
				Messages: messages,
				Users:    users,
			},
		},
	}, nil
}

func searchID[T any, PT interface {
	GetId() string
	*T
}](id string, list []PT) bool {
	for _, v := range list {
		if v.GetId() == id {
			return true
		}
	}

	return false
}

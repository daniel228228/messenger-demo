package service

import (
	"github.com/pkg/errors"
	"messenger.api/go/api"
)

func (s *service) GetDialog(userID string, req *api.GetDialogRequest) (*api.GetDialogResponse, error) {
	var peerID string

	switch v := req.Peer.Peer.(type) {
	case *api.Peer_User:
		peerID = v.User.UserId
	default:
		return nil, errors.New("unsupported peer type")
	}

	lastMsg, lastMsgID := s.getLastMessage(userID, peerID)
	if lastMsg == nil {
		return nil, ErrDialogNotFound
	}

	lastReadMsgIDFromUser, _ := s.getLastReadMessageID(userID, peerID)
	lastReadMsgIDFromPeer, _ := s.getLastReadMessageID(peerID, userID)
	countUnreadMsgs, _ := s.getCountUnreadMessages(userID, peerID, lastReadMsgIDFromUser)

	return &api.GetDialogResponse{
		Dialog:      s.buildDialogProto(req.Peer, lastMsgID, countUnreadMsgs, lastReadMsgIDFromUser, lastReadMsgIDFromPeer),
		LastMessage: lastMsg,
		User:        cachedUsers.get(peerID),
	}, nil
}

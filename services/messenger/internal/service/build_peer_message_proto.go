package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"messenger.api/go/api"

	"messenger.messenger/internal/models/entities"
)

func (s *service) buildPeerMessageProto(message *entities.Message, outgoing bool, peer, author *api.Peer) *api.Message {
	if message == nil {
		return nil
	}

	return &api.Message{
		Message: &api.Message_PeerMessage{
			PeerMessage: &api.PeerMessage{
				Id:       message.ID.String(),
				Date:     timestamppb.New(message.Timestamp),
				Message:  message.Message,
				Peer:     peer,
				Outgoing: outgoing,
				FromPeer: author,
			},
		},
	}
}

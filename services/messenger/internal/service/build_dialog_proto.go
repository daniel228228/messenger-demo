package service

import "messenger.api/go/api"

func (s *service) buildDialogProto(peer *api.Peer, lastMsgID string, unreadCount int, lastInboxReadMsgID, lastOutboxReadMessageID *string) *api.Dialog {
	return &api.Dialog{
		Peer:                    peer,
		LastMessageId:           lastMsgID,
		UnreadCount:             int32(unreadCount),
		ReadInboxLastMessageId:  lastInboxReadMsgID,
		ReadOutboxLastMessageId: lastOutboxReadMessageID,
	}
}

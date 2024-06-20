package entities

import (
	"github.com/google/uuid"
)

type MessagePeerUserToUser struct {
	MessageID  uuid.UUID `db:"message_id"`
	Message    *Message
	FromUserID uuid.UUID `db:"from_user_id"`
	ToUserID   uuid.UUID `db:"to_user_id"`
}

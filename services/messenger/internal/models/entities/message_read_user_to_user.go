package entities

import (
	"github.com/google/uuid"
)

type MessageReadUserToUser struct {
	FromUserID        uuid.UUID `db:"from_user_id"`
	ToUserID          uuid.UUID `db:"to_user_id"`
	LastReadMessageID uuid.UUID `db:"last_read_message_id"`
	LastReadMessage   *Message
}

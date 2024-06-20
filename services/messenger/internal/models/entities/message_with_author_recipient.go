package entities

import (
	"github.com/google/uuid"
)

type MessageWithAuthorRecipient struct {
	Message
	Author    uuid.UUID `db:"from_user_id"`
	Recipient uuid.UUID `db:"to_user_id"`
}

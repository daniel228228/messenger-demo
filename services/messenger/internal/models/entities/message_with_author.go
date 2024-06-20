package entities

import (
	"github.com/google/uuid"
)

type MessageWithAuthor struct {
	Message
	Author uuid.UUID `db:"from_user_id"`
}

package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `db:"id"`
	Timestamp time.Time `db:"timestamp"`
	Message   string    `db:"message"`
}

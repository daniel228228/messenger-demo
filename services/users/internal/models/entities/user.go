package entities

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	FirstName *string   `db:"first_name"`
	LastName  *string   `db:"last_name"`
}

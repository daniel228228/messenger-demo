package repo

import (
	"github.com/google/uuid"
	"messenger.users/internal/models/entities"
)

func (r *repo) CreateUser(user *entities.User) (uuid.UUID, error) {
	query1 := `SELECT id FROM users WHERE users.username = $1`
	query2 := `INSERT INTO users (username, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`

	var id uuid.UUID

	err := r.db.Select(&id, query1, user.Username)
	if err != nil {
		return id, err
	}

	if id != uuid.Nil {
		return id, nil
	}

	rows, err := r.db.Query(query2, user.Username, user.FirstName, user.LastName)
	if rows.Next() {
		rows.Scan(&id)
	}

	return id, err
}

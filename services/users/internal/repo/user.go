package repo

import (
	"database/sql"
	"errors"

	"messenger.users/internal/models/entities"
)

func (r *repo) User(userID string) (*entities.User, error) {
	query := `SELECT * FROM user WHERE id = $1`

	user := &entities.User{}

	if err := r.db.Select(user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

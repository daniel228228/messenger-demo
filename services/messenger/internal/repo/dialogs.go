package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"messenger.messenger/internal/models/entities"
)

func (r *repo) Dialogs(userID string, offsetTimestamp time.Time, limit int, messages *[]entities.MessageWithAuthorRecipient) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, ErrBadUserID
	}

	if err := r.db.Select(messages, dialogQueryOffset, userID, offsetTimestamp.Truncate(time.Microsecond), limit); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	var count int64

	rows, err := r.db.Query(dialogQueryTotal, userUUID)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		rows.Scan(&count)
	}

	return int(count), nil
}

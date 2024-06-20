package repo

import (
	"time"

	"github.com/google/uuid"

	"messenger.messenger/internal/models/entities"
)

func (r *repo) Messages(userID, peerID string, offsetID *string, limit int, messages *[]entities.MessageWithAuthor) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return 0, ErrBadPeerID
	}

	if offsetID == nil {
		err := r.db.Select(messages, msgQuery, userUUID, peerUUID, limit)
		if err != nil {
			return 0, err
		}
	} else {
		offsetUUID, err := uuid.Parse(*offsetID)
		if err != nil {
			return 0, ErrBadMessageID
		}

		var timestamp time.Time

		rows, err := r.db.Query(msgQueryTimestamp, userUUID, peerUUID, offsetUUID)
		if err != nil {
			return 0, err
		}
		if rows.Next() {
			rows.Scan(&timestamp)
		}

		if timestamp.IsZero() {
			return 0, ErrBadMessageID
		}

		if err := r.db.Select(messages, msgQueryOffset, timestamp.Truncate(time.Microsecond), offsetUUID, userUUID, peerUUID, limit); err != nil {
			return 0, err
		}
	}

	var count int64

	rows, err := r.db.Query(msgQueryTotal, userUUID, peerUUID)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		rows.Scan(&count)
	}

	return int(count), nil
}

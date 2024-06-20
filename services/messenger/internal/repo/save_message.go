package repo

import (
	"github.com/google/uuid"

	"messenger.messenger/internal/models/entities"
)

func (r *repo) SaveMessage(userID, peerID string, message *entities.Message) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return ErrBadPeerID
	}

	query := `WITH msg AS (
		INSERT INTO message(message)
			VALUES ($1)
			RETURNING id
		)
		INSERT INTO message_peer_user_to_user (from_user_id, message_id, to_user_id)
		SELECT $2, msg.id, $3
		FROM msg
	`

	_, err = r.db.Exec(query, message.Message, userUUID, peerUUID)
	if err != nil {
		return err
	}

	return nil
}

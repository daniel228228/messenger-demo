package repo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"messenger.messenger/internal/models/entities"
)

func (r *repo) LastReadMessageID(userID, peerID string) (*string, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return nil, ErrBadPeerID
	}

	readMsg := &entities.MessageReadUserToUser{}

	query := `SELECT last_read_message_id FROM message_read_user_to_user WHERE from_user_id = $1 AND to_user_id = $2`

	if err := r.db.Select(readMsg, query, userUUID, peerUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	resp := readMsg.LastReadMessageID.String()

	return &resp, nil
}

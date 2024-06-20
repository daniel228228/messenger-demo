package repo

import (
	"github.com/google/uuid"
)

func (r *repo) ReadMessages(userID, peerID string, lastID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return ErrBadPeerID
	}

	lastUUID, err := uuid.Parse(lastID)
	if err != nil {
		return ErrBadMessageID
	}

	res, err := r.db.Exec(msgReadQuery, userUUID, peerUUID, lastUUID)
	if err != nil {
		return err
	}
	if v, err := res.RowsAffected(); err != nil && v == 0 {
		return ErrReadRejected
	}

	return nil
}

package repo

import (
	"github.com/google/uuid"
)

func (r *repo) IsUnreadDialog(userID, peerID string) (bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return false, ErrBadPeerID
	}

	res, err := r.db.Exec(isUnreadDialogQuery, userUUID, peerUUID)
	if err != nil {
		return false, err
	}
	if v, err := res.RowsAffected(); err != nil && v == 0 {
		return false, nil
	}

	return true, nil
}

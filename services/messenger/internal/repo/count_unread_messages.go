package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"messenger.messenger/internal/models/entities"
)

func (r *repo) CountUnreadMessages(userID, peerID string, lastReadMsgID *string) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, ErrBadUserID
	}

	peerUUID, err := uuid.Parse(peerID)
	if err != nil {
		return 0, ErrBadPeerID
	}

	var count int64

	if lastReadMsgID == nil {
		rows, err := r.db.Query(msgUnreadQueryCount, peerUUID, userUUID)
		if err != nil {
			return 0, err
		}
		if rows.Next() {
			rows.Scan(&count)
		}
	} else {
		lastReadMsgUUID, err := uuid.Parse(*lastReadMsgID)
		if err != nil {
			return 0, ErrBadMessageID
		}

		message := &entities.Message{}

		query := `SELECT timestamp FROM message WHERE id = $1`

		if err := r.db.Select(message, query, lastReadMsgUUID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, ErrBadMessageID
			}

			return 0, err
		}

		rows, err := r.db.Query(msgUnreadQueryCountOffset, message.Timestamp.Truncate(time.Microsecond), message.ID, peerUUID, userUUID)
		if err != nil {
			return 0, err
		}
		if rows.Next() {
			rows.Scan(&count)
		}
	}

	return int(count), nil
}

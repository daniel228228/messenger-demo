package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"messenger.messenger/pkg/log"

	"messenger.messenger/internal/models/entities"
	"messenger.messenger/pkg/config"
)

var Name = "Repo"

type Repo interface {
	Dialogs(userID string, offsetTimestamp time.Time, limit int, messages *[]entities.MessageWithAuthorRecipient) (int, error)
	Messages(userID, peerID string, offsetID *string, limit int, messages *[]entities.MessageWithAuthor) (int, error)
	SaveMessage(userID, peerID string, message *entities.Message) error
	ReadMessages(userID, peerID string, lastID string) error
	LastReadMessageID(userID, peerID string) (*string, error)
	CountUnreadMessages(userID, peerID string, lastReadMsgID *string) (int, error)
	IsUnreadDialog(userID, peerID string) (bool, error)
}

const (
	ConnRetries = 20
)

var (
	ErrBadUserID    = errors.New("bad user id")
	ErrBadPeerID    = errors.New("bad peer id")
	ErrBadMessageID = errors.New("bad message id")
	ErrPeerNotFound = errors.New("peer not found")
	ErrReadRejected = errors.New("read rejected")
)

type repo struct {
	config *config.Config
	log    *log.Logger
	db     *sqlx.DB
}

func NewRepo(config *config.Config, log *log.Logger) *repo {
	return &repo{
		config: config,
		log:    log,
	}
}

func (r *repo) Start(ctx context.Context) error {
	period := 16 * time.Millisecond

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	var err error

	for i := 1; i <= ConnRetries; i++ {
		select {
		case <-ctx.Done():
			return err
		case <-ticker.C:
		}

		r.log.Trace().Msgf("Attempt %d: connecting to DB", i)

		err = r.connect()
		if err != nil {
			r.log.Trace().Msgf("Attempt %d failed: DB error: %v", i, err)

			if period < 5*time.Second {
				period *= 2
			} else {
				period = 10 * time.Second
			}

			ticker.Reset(period)
		} else {
			r.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()
			return nil
		}
	}

	r.log.Error().Err(err).Msg("Can't connect to DB")
	return err
}

func (r *repo) connect() error {
	db, err := sqlx.Connect("postgres", r.config.DBUrl)
	if err != nil {
		return err
	}

	r.db = db
	return nil
}

func (r *repo) Stop(ctx context.Context) error {
	r.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return r.db.DB.Close()
}

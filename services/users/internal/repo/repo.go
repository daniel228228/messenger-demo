package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"messenger.users/pkg/log"

	"messenger.users/internal/models/entities"
	"messenger.users/pkg/config"
)

var Name = "Repo"

type Repo interface {
	User(userID string) (*entities.User, error)
	CreateUser(user *entities.User) (uuid.UUID, error)
}

const (
	ConnRetries = 20
)

var (
	ErrUserNotFound = errors.New("user not found")
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

package postresrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type Logger struct {
	db *pgxpool.Pool
}

func NewLogger(db *pgxpool.Pool) *Logger {
	return &Logger{db}
}

func (l *Logger) Save(ctx context.Context, log *model.PlaybackLog) (*model.PlaybackLog, error) {
	var id int
	err := l.db.QueryRow(
		ctx,
		`INSERT INTO playback_logs (track_id, played_at, amount_paid) 
			VALUES ($1, $2, $3)
		RETURNING track_id`,
		log.TrackID().Int(), log.PlayedAt(), log.AmountPaid(),
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				err = repo.ErrNonExistingData
			}
		}
		return nil, err
	}

	log, err = model.RebuildPlaybackLog(id, log.TrackID().Int(), log.PlayedAt(), log.AmountPaid())
	if err != nil {
		return nil, fmt.Errorf("failed to rebuild playback log: %w", err)
	}

	return log, nil
}

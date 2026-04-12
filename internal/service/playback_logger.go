package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type logRepo interface {
	Save(ctx context.Context, log *model.PlaybackLog) (*model.PlaybackLog, error)
}

type PlaybackLogger struct {
	repo logRepo
}

func NewPlaybackLogger(repo logRepo) *PlaybackLogger {
	return &PlaybackLogger{repo}
}

func (l *PlaybackLogger) RegisterPlayback(ctx context.Context, trackID, amountPaid int) (*model.PlaybackLog, error) {
	trackIdParsed, err := model.ParseTrackID(trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse track id: %w", ErrInvalidTrackID)
	}

	log, err := model.NewPlayBackLog(trackIdParsed, amountPaid)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaidAmount) {
			err = ErrInvalidPaidAmount
		}

		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	log, err = l.repo.Save(ctx, log)
	if err != nil {
		if errors.Is(err, repo.ErrNonExistingData) {
			err = ErrNonExistingTrack
		}
		return nil, fmt.Errorf("faild to save playback log: %w", err)
	}

	return log, nil
}

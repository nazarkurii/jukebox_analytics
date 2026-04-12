package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type trackRepo interface {
	GetTrack(ctx context.Context, id model.TrackID) (*model.Track, error)
	UpdateTrack(ctx context.Context, track *model.Track) error
}

type Track struct {
	repo trackRepo
}

func NewTrack(repo trackRepo) *Track {
	return &Track{repo}
}

func (t *Track) UpdatePrice(ctx context.Context, id, newPrice int) (*model.Track, error) {
	trackID, err := model.ParseTrackID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id: %w", ErrInvalidTrackID)
	}

	track, err := t.repo.GetTrack(ctx, trackID)
	if err != nil {
		if errors.Is(err, repo.ErrNonExistingData) {
			err = ErrNonExistingTrack
		}
		return nil, fmt.Errorf("failed to retrieve track: %w", err)
	}

	err = track.UpdatePrice(newPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to update track price: %w", ErrInvalidTrackPrice)
	}

	err = t.repo.UpdateTrack(ctx, track)
	if err != nil {
		if errors.Is(err, repo.ErrNonExistingData) {
			err = ErrNonExistingTrack
		}
		return nil, fmt.Errorf("failed to update track: %w", err)
	}

	return track, nil
}

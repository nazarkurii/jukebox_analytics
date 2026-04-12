package service

import (
	"context"
	"fmt"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
)

type statsRepo interface {
	GetPopularTracks(ctx context.Context, number int) ([]*model.Track, error)
	GetRevenue(ctx context.Context) ([]struct {
		ID           int
		Name, Artist string
		Revenue      int
	}, error)
}

type Stats struct {
	repo statsRepo
}

func NewStats(repo statsRepo) *Stats {
	return &Stats{repo}
}

func (s *Stats) GetTopTracks(ctx context.Context) ([]*model.Track, error) {
	tracks, err := s.repo.GetPopularTracks(ctx, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve popular tracks: %w", err)
	}

	return tracks, nil
}

func (s *Stats) GetRevenue(ctx context.Context) ([]struct {
	ID           int
	Name, Artist string
	Revenue      int
}, int, error) {
	trackRevenues, err := s.repo.GetRevenue(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve revenue data: %w", err)
	}

	var total int
	for _, trackRevenue := range trackRevenues {
		total += trackRevenue.Revenue
	}

	return trackRevenues, total, nil
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/model/testmodel"
)

type statsRepoMock struct {
	t *testing.T

	expectedNumber int

	tracks          []*model.Track
	returnTracksErr error
	revenues        []struct {
		ID           int
		Name, Artist string
		Revenue      int
	}
	returnRevenueErr error
}

func (m *statsRepoMock) init(t *testing.T) {
	m.t = t
}

func (m *statsRepoMock) GetPopularTracks(ctx context.Context, number int) ([]*model.Track, error) {
	if number != m.expectedNumber {
		m.t.Fatalf("expected number, passed into GetPopularTracks(ctx, number), to be %q, got %q", m.expectedNumber, number)
	}
	return m.tracks, m.returnTracksErr
}

func (m *statsRepoMock) GetRevenue(ctx context.Context) ([]struct {
	ID           int
	Name, Artist string
	Revenue      int
}, error) {
	return m.revenues, m.returnRevenueErr
}

func TestStats_GetTopTracks(t *testing.T) {
	t.Parallel()

	unexpectedErr := errors.New("connection interrupted")

	tests := []struct {
		name          string
		statsRepoMock statsRepoMock
		expectedErr   error
	}{
		{
			name: "returns error when repo.GetPopularTracks fails",
			statsRepoMock: statsRepoMock{
				expectedNumber:  3,
				returnTracksErr: unexpectedErr,
			},
			expectedErr: unexpectedErr,
		},
		{
			name: "returns tracks when succeeds",
			statsRepoMock: statsRepoMock{
				expectedNumber: 3,
				tracks: []*model.Track{
					testmodel.TestTrack(t),
					testmodel.TestTrack(t),
					testmodel.TestTrack(t),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.statsRepoMock.init(t)

			statsService := NewStats(&tc.statsRepoMock)
			tracks, err := statsService.GetTopTracks(context.Background())

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %q, got %q", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if tracks == nil || len(tracks) != len(tc.statsRepoMock.tracks) {
					t.Errorf("expected tracks length to be %d, got %d", len(tc.statsRepoMock.tracks), len(tracks))
				}
			} else if tracks != nil {
				t.Errorf("expected tracks to be nil, got %v", tracks)
			}
		})
	}
}

func TestStats_GetRevenue(t *testing.T) {
	t.Parallel()

	unexpectedErr := errors.New("connection interrupted")

	tests := []struct {
		name          string
		statsRepoMock statsRepoMock
		expectedTotal int
		expectedErr   error
	}{
		{
			name: "returns error when repo.GetRevenue fails",
			statsRepoMock: statsRepoMock{
				returnRevenueErr: unexpectedErr,
			},
			expectedErr: unexpectedErr,
		},
		{
			name: "returns revenues and total when succeeds",
			statsRepoMock: statsRepoMock{
				revenues: []struct {
					ID           int
					Name, Artist string
					Revenue      int
				}{
					{ID: 1, Name: "A", Artist: "X", Revenue: 100},
					{ID: 2, Name: "B", Artist: "Y", Revenue: 200},
					{ID: 3, Name: "C", Artist: "Z", Revenue: 300},
				},
			},
			expectedTotal: 600,
		},
		{
			name: "returns zero total when no revenues",
			statsRepoMock: statsRepoMock{
				revenues: []struct {
					ID           int
					Name, Artist string
					Revenue      int
				}{},
			},
			expectedTotal: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.statsRepoMock.init(t)

			statsService := NewStats(&tc.statsRepoMock)
			revenues, total, err := statsService.GetRevenue(context.Background())

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %q, got %q", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if total != tc.expectedTotal {
					t.Errorf("expected total to be %d, got %d", tc.expectedTotal, total)
				}

				if len(revenues) != len(tc.statsRepoMock.revenues) {
					t.Errorf("expected revenues length to be %d, got %d", len(tc.statsRepoMock.revenues), len(revenues))
				}
			} else {
				if revenues != nil {
					t.Errorf("expected revenues to be nil, got %v", revenues)
				}
				if total != 0 {
					t.Errorf("expected total to be 0, got %d", total)
				}
			}
		})
	}
}

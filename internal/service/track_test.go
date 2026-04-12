package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/model/testmodel"
	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type trackRepoMock struct {
	t               *testing.T
	expectedID      model.TrackID
	track           *model.Track
	returnGetErr    error
	returnUpdateErr error
}

func (m *trackRepoMock) init(t *testing.T) {
	m.t = t
}

func (m *trackRepoMock) GetTrack(ctx context.Context, id model.TrackID) (*model.Track, error) {
	if id != m.expectedID {
		m.t.Fatalf("expected id, passed into GetTrack(ctx, id), to be %q, got %q", m.expectedID, id)
	}
	return m.track, m.returnGetErr
}

func (m *trackRepoMock) UpdateTrack(ctx context.Context, track *model.Track) error {
	if track != m.track {
		m.t.Fatalf("expected track, passed into GetTrack(ctx, id), to be '%v', got '%v'", m.track, track)
	}
	return m.returnUpdateErr
}

func TestTrack_UpdatePrice(t *testing.T) {
	t.Parallel()

	unexpectedErr := errors.New("connection interupted")

	tests := []struct {
		name          string
		id            int
		newPrice      int
		trackRepoMock trackRepoMock
		expectedErr   error
	}{
		{
			name:        "returns error when id is not valid",
			id:          testmodel.TestInvalidID(),
			expectedErr: ErrInvalidTrackID,
		},
		{
			name: "returns error when track does not exist",
			id:   testmodel.TestTrackID(t).Int(),
			trackRepoMock: trackRepoMock{
				expectedID:   testmodel.TestTrackID(t),
				returnGetErr: repo.ErrNonExistingData,
			},
			expectedErr: ErrNonExistingTrack,
		},
		{
			name: "returns error when repo.GetTrack unexpectedly fails",
			id:   testmodel.TestTrackID(t).Int(),
			trackRepoMock: trackRepoMock{
				expectedID:   testmodel.TestTrackID(t),
				returnGetErr: unexpectedErr,
			},
			expectedErr: unexpectedErr,
		},
		{
			name:     "returns error when newPrice is not valid",
			id:       testmodel.TestTrackID(t).Int(),
			newPrice: testmodel.TestInvalidTrackPrice(),
			trackRepoMock: trackRepoMock{
				expectedID: testmodel.TestTrackID(t),
				track:      testmodel.TestTrack(t),
			},
			expectedErr: ErrInvalidTrackPrice,
		},

		{
			name:     "returns error when repo.UpdateTrack fails due to non-existing track (deleted betweed get and update calls)",
			id:       testmodel.TestTrackID(t).Int(),
			newPrice: testmodel.TestValidTrackPrice(),
			trackRepoMock: trackRepoMock{
				expectedID:      testmodel.TestTrackID(t),
				track:           testmodel.TestTrack(t),
				returnUpdateErr: repo.ErrNonExistingData,
			},
			expectedErr: ErrNonExistingTrack,
		},

		{
			name:     "returns error when repo.UpdateTrack fails unexpectedly",
			id:       testmodel.TestTrackID(t).Int(),
			newPrice: testmodel.TestValidTrackPrice(),
			trackRepoMock: trackRepoMock{
				expectedID:      testmodel.TestTrackID(t),
				track:           testmodel.TestTrack(t),
				returnUpdateErr: unexpectedErr,
			},
			expectedErr: unexpectedErr,
		},

		{
			name:     "returns track with updated price when succeeds",
			id:       testmodel.TestTrackID(t).Int(),
			newPrice: testmodel.TestValidTrackPrice(),
			trackRepoMock: trackRepoMock{
				expectedID: testmodel.TestTrackID(t),
				track:      testmodel.TestTrack(t),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.trackRepoMock.init(t)

			trackService := NewTrack(&tc.trackRepoMock)
			track, err := trackService.UpdatePrice(context.Background(), tc.id, tc.newPrice)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %q, got %q", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if track != nil && track.Price() != tc.newPrice {
					t.Errorf("expected track price to be %q, got %q", tc.newPrice, track.Price())
				}
			} else if track != nil {
				t.Errorf("expected track to be nil, got %q", track.Price())
			}

		})
	}
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/model/testmodel"

	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type logRepoMock struct {
	t         *testing.T
	returnErr error
}

func (m *logRepoMock) init(t *testing.T) {
	m.t = t
}

func (m *logRepoMock) Save(ctx context.Context, log *model.PlaybackLog) (*model.PlaybackLog, error) {
	if log == nil {
		m.t.Fatalf("expected log, passed into Save(ctx, log) not to be nil, got nil")
	}

	if m.returnErr != nil {
		return nil, m.returnErr
	}

	log = testmodel.TestPlaybackLogWithID(m.t, log)

	return log, m.returnErr
}

func TestPlaybackLogger_RegisterPlayback(t *testing.T) {
	t.Parallel()

	unexpectedErr := errors.New("connection interrupted")

	tests := []struct {
		name        string
		trackID     int
		amountPaid  int
		logRepoMock logRepoMock
		expectedErr error
	}{
		{
			name:        "returns error when trackID is invalid",
			trackID:     testmodel.TestInvalidID(),
			amountPaid:  testmodel.TestValidAmountPaid(),
			expectedErr: ErrInvalidTrackID,
		},
		{
			name:        "returns error when amountPaid is invalid",
			trackID:     testmodel.TestTrackID(t).Int(),
			amountPaid:  testmodel.TestInvalidAmountPaid(),
			expectedErr: ErrInvalidPaidAmount,
		},
		{
			name:       "returns error when repo.Save fails due to non-existing track",
			trackID:    testmodel.TestTrackID(t).Int(),
			amountPaid: testmodel.TestValidAmountPaid(),
			logRepoMock: logRepoMock{
				returnErr: repo.ErrNonExistingData,
			},
			expectedErr: ErrNonExistingTrack,
		},
		{
			name:       "returns error when repo.Save fails unexpectedly",
			trackID:    testmodel.TestTrackID(t).Int(),
			amountPaid: testmodel.TestValidAmountPaid(),
			logRepoMock: logRepoMock{
				returnErr: unexpectedErr,
			},
			expectedErr: unexpectedErr,
		},
		{
			name:        "returns playback log when succeeds",
			trackID:     testmodel.TestTrackID(t).Int(),
			amountPaid:  testmodel.TestValidAmountPaid(),
			logRepoMock: logRepoMock{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.logRepoMock.init(t)

			logger := NewPlaybackLogger(&tc.logRepoMock)
			log, err := logger.RegisterPlayback(context.Background(), tc.trackID, tc.amountPaid)

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %q, got %q", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if log == nil {
					t.Errorf("expected log to be non-nil")
				}
			} else if log != nil {
				t.Errorf("expected log to be nil, got %v", log)
			}
		})
	}
}

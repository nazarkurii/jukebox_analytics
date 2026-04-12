package model

import (
	"errors"
	"testing"
	"time"
)

func TestParsePlaybackLogID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       int
		expectedID  int
		expectedErr error
	}{
		{
			name:        "returns error when id is zero",
			input:       0,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "returns error when id is negative",
			input:       -1,
			expectedErr: ErrInvalidID,
		},
		{
			name:       "returns PlaybackLogID when id is positive number",
			input:      1,
			expectedID: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := ParsePlaybackLogID(tc.input)

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if id.id != tc.expectedID {
					t.Errorf("expected id to be %d, got %d", tc.expectedID, id.id)
				}
			} else {
				if id.id != 0 {
					t.Errorf("expected id to be zero value, got %d", id.id)
				}
			}
		})
	}
}

func TestPlaybackLogID_Int(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   PlaybackLogID
	}{
		{
			name: "returns correct int value",
			id:   PlaybackLogID{id: 1},
		},
		{
			name: "returns correct int value",
			id:   PlaybackLogID{id: 999},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if result := tc.id.Int(); result != tc.id.id {
				t.Errorf("expected %d, got %d", tc.id.id, result)
			}
		})
	}
}

func TestPlaybackLog_Getters(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	id, _ := ParsePlaybackLogID(101)
	trackID, _ := ParseTrackID(202)

	tests := []struct {
		name string
		log  PlaybackLog
	}{
		{
			name: "returns correct values from log fields",
			log: PlaybackLog{
				id:         id,
				trackID:    trackID,
				playedAt:   now,
				amountPaid: 50,
			},
		},
		{
			name: "returns correct values from log fields",
			log: PlaybackLog{
				id:         PlaybackLogID{id: 1},
				trackID:    TrackID{id: 2},
				playedAt:   now.Add(time.Second),
				amountPaid: 250,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if result := tc.log.ID(); result != tc.log.id {
				t.Errorf("expected id to be %v, got %v", tc.log.id, result)
			}

			if result := tc.log.TrackID(); result != tc.log.trackID {
				t.Errorf("expected track id to be %v, got %v", tc.log.trackID, result)
			}

			if result := tc.log.PlayedAt(); !result.Equal(tc.log.playedAt) {
				t.Errorf("expected played at to be %v, got %v", tc.log.playedAt, result)
			}

			if result := tc.log.AmountPaid(); result != tc.log.amountPaid {
				t.Errorf("expected amount paid to be %d, got %d", tc.log.amountPaid, result)
			}
		})
	}
}

func TestNewPlayBackLog(t *testing.T) {
	t.Parallel()

	validTrackID, _ := ParseTrackID(1)

	tests := []struct {
		name        string
		trackID     TrackID
		amountPaid  int
		expectedErr error
	}{
		{
			name:        "returns error when amountPaid is invalid",
			trackID:     validTrackID,
			amountPaid:  0,
			expectedErr: ErrInvalidPaidAmount,
		},
		{
			name:        "returns error when amountPaid is invalid",
			trackID:     validTrackID,
			amountPaid:  -1,
			expectedErr: ErrInvalidPaidAmount,
		},
		{
			name:       "returns playback log when valid",
			trackID:    validTrackID,
			amountPaid: 100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log, err := NewPlayBackLog(tc.trackID, tc.amountPaid)

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected error to be %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedErr == nil {
				if log == nil {
					t.Fatal("expected playback log not to be nil")
				}

				if log.id.id != 0 {
					t.Errorf("expected id to be %d, got %d", 0, log.id.id)
				}

				if log.trackID != tc.trackID {
					t.Errorf("expected track id to be %v, got %v", tc.trackID, log.trackID)
				}

				if log.amountPaid != tc.amountPaid {
					t.Errorf("expected amount paid to be %d, got %d", tc.amountPaid, log.amountPaid)
				}

				if time.Since(log.playedAt) > time.Second {
					t.Error("expected played_at to be now")
				}
			} else {
				if log != nil {
					t.Errorf("expected playback log to be nil, got %v", log)
				}
			}
		})
	}
}

func TestRebuildPlaybackLog(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	tests := []struct {
		name        string
		id          int
		trackID     int
		playedAt    time.Time
		amountPaid  int
		expectedErr error
	}{
		{
			name:        "returns error when id is invalid",
			id:          0,
			trackID:     1,
			playedAt:    now,
			amountPaid:  100,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "returns error when track id is invalid",
			id:          1,
			trackID:     0,
			playedAt:    now,
			amountPaid:  100,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "returns error when amountPaid is invalid",
			id:          1,
			trackID:     1,
			playedAt:    now,
			amountPaid:  0,
			expectedErr: ErrInvalidPaidAmount,
		},
		{
			name:        "returns error when playedAt is zero",
			id:          1,
			trackID:     1,
			playedAt:    time.Time{},
			amountPaid:  100,
			expectedErr: errors.New("invalid playedAt value"),
		},
		{
			name:       "returns rebuilt log when all fields are valid",
			id:         10,
			trackID:    20,
			playedAt:   now,
			amountPaid: 150,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log, err := RebuildPlaybackLog(tc.id, tc.trackID, tc.playedAt, tc.amountPaid)

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error to be %q, got %q", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			if log == nil {
				t.Fatal("expected playback log not to be nil")
			}

			if log.id.id != tc.id {
				t.Errorf("expected id to be %d, got %d", tc.id, log.id.id)
			}

			if log.trackID.id != tc.trackID {
				t.Errorf("expected trackID to be %d, got %d", tc.trackID, log.trackID.id)
			}

			if !log.playedAt.Equal(tc.playedAt) {
				t.Errorf("expected playedAt to be %v, got %v", tc.playedAt, log.playedAt)
			}

			if log.amountPaid != tc.amountPaid {
				t.Errorf("expected amountPaid to be %d, got %d", tc.amountPaid, log.amountPaid)
			}
		})
	}
}

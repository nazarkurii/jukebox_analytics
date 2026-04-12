package model

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseTrackID(t *testing.T) {
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
			name:       "returns TrackID when id is positive number",
			input:      1,
			expectedID: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := ParseTrackID(tc.input)

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

func TestTrackID_Int(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   TrackID
	}{
		{
			name: "returns correct int value",
			id:   TrackID{id: 10},
		},
		{
			name: "returns correct value for another id",
			id:   TrackID{id: 999},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if result := tc.id.Int(); result != tc.id.id {
				t.Errorf("expected id to be %d, got %d", tc.id.id, result)
			}

			if result := tc.id.String(); result != strconv.Itoa(tc.id.id) {
				t.Errorf("expected id to be %s, got %s", strconv.Itoa(tc.id.id), result)
			}
		})
	}
}

func TestTrack_Getters(t *testing.T) {
	t.Parallel()

	id1, _ := ParseTrackID(500)
	id2, _ := ParseTrackID(1)

	tests := []struct {
		name  string
		track *Track
	}{
		{
			name: "returns correct values from track fields",
			track: &Track{
				id:     id1,
				name:   "Midnight City",
				artist: "M83",
				price:  120,
			},
		},
		{
			name: "returns correct values from track fields",
			track: &Track{
				id:     id2,
				name:   "John",
				artist: "Doe",
				price:  380,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if result := tc.track.ID(); result != tc.track.id {
				t.Errorf("expected id to be %v, got %v", tc.track.id, result)
			}

			if result := tc.track.Name(); result != tc.track.name {
				t.Errorf("expected name to be %q, got %q", tc.track.name, result)
			}

			if result := tc.track.Artist(); result != tc.track.artist {
				t.Errorf("expected artist to be %q, got %q", tc.track.artist, result)
			}

			if result := tc.track.Price(); result != tc.track.price {
				t.Errorf("expected price to be %d, got %d", tc.track.price, result)
			}

		})
	}
}

func TestTrack_UpdatePrice(t *testing.T) {
	t.Parallel()

	id, _ := ParseTrackID(1)

	tests := []struct {
		name          string
		initialTrack  *Track
		newPrice      int
		expectedPrice int
		expectedErr   error
	}{
		{
			name:          "returns error when price is zero",
			initialTrack:  &Track{id: id, price: 100},
			newPrice:      0,
			expectedPrice: 100,
			expectedErr:   ErrInvalidTrackPrice,
		},
		{
			name:          "returns error when price is negative",
			initialTrack:  &Track{id: id, price: 100},
			newPrice:      -50,
			expectedPrice: 100,
			expectedErr:   ErrInvalidTrackPrice,
		},
		{
			name:          "updates price when value is valid",
			initialTrack:  &Track{id: id, price: 100},
			newPrice:      200,
			expectedPrice: 200,
			expectedErr:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.initialTrack.UpdatePrice(tc.newPrice)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error to be %v, got %v", tc.expectedErr, err)
			}

			if tc.initialTrack.price != tc.expectedPrice {
				t.Errorf("expected price to be %d, got %d", tc.expectedPrice, tc.initialTrack.price)
			}
		})
	}
}

func TestRebuildTrack(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          int
		trackName   string
		artist      string
		price       int
		expectedErr error
	}{
		{
			name:        "returns error when id is invalid",
			id:          0,
			trackName:   "Song A",
			artist:      "Artist A",
			price:       100,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "returns error when name is empty",
			id:          1,
			trackName:   "",
			artist:      "Artist A",
			price:       100,
			expectedErr: errors.New("invalid track name (empty)"),
		},
		{
			name:        "returns error when artist is empty",
			id:          1,
			trackName:   "Song A",
			artist:      "",
			price:       100,
			expectedErr: errors.New("invalid track artist (empty)"),
		},
		{
			name:        "returns error when price is invalid",
			id:          1,
			trackName:   "Song A",
			artist:      "Artist A",
			price:       0,
			expectedErr: ErrInvalidTrackPrice,
		},
		{
			name:        "returns error when price is invalid",
			id:          1,
			trackName:   "Song A",
			artist:      "Artist A",
			price:       -1,
			expectedErr: ErrInvalidTrackPrice,
		},
		{
			name:      "returns rebuilt track when all args are valid",
			id:        1,
			trackName: "Song A",
			artist:    "Artist A",
			price:     100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			track, err := RebuildTrack(tc.id, tc.trackName, tc.artist, tc.price)

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error msg to be %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if track == nil {
				t.Fatal("expected track not to be nil")
			}

			if track.id.id != tc.id {
				t.Errorf("expected id to be %d, got %d", tc.id, track.id.Int())
			}

			if track.name != tc.trackName {
				t.Errorf("expected name to be %q, got %q", tc.trackName, track.name)
			}

			if track.artist != tc.artist {
				t.Errorf("expected artist to be %q, got %q", tc.artist, track.artist)
			}

			if track.price != tc.price {
				t.Errorf("expected price to be %d, got %d", tc.price, track.price)
			}
		})
	}
}

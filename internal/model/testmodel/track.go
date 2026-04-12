package testmodel

import (
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
)

func TestTrackID(t *testing.T) model.TrackID {
	id, err := model.ParseTrackID(1)
	if err != nil {
		t.Fatalf("failed to create test track id: %v", err)
	}

	return id
}

func TestInvalidID() int {
	return -1
}

func TestValidTrackPrice() int {
	return 100
}

func TestInvalidTrackPrice() int {
	return 0
}

func TestTrack(t *testing.T) *model.Track {
	track, err := model.RebuildTrack(TestTrackID(t).Int(), "track", "artist", TestValidTrackPrice())
	if err != nil {
		t.Fatalf("failed to create test track: %v", err)
	}

	return track
}

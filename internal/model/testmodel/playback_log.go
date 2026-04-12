package testmodel

import (
	"log"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
)

func TestPlaybackLogID() model.PlaybackLogID {
	id, err := model.ParsePlaybackLogID(1)
	if err != nil {
		log.Fatalf("failed to create test playback log id: %v", err)
	}

	return id
}
func TestValidAmountPaid() int {
	return 100
}

func TestInvalidAmountPaid() int {
	return 0
}

func TestPlaybackLogWithID(t *testing.T, playbackLog *model.PlaybackLog) *model.PlaybackLog {
	playbackLog, err := model.RebuildPlaybackLog(1, playbackLog.TrackID().Int(), playbackLog.PlayedAt(), playbackLog.AmountPaid())
	if err != nil {
		t.Fatalf("failed to create test playback log id: %s", err)
	}
	return playbackLog
}

func TestPlaybackLog(t *testing.T) *model.PlaybackLog {
	log, err := model.NewPlayBackLog(TestTrackID(t), TestValidAmountPaid())
	if err != nil {
		t.Fatalf("failed to create test playback log: %s", err)
	}

	return TestPlaybackLogWithID(t, log)
}

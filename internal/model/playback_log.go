package model

import (
	"errors"
	"time"
)

type PlaybackLogID struct {
	id int
}

func (id PlaybackLogID) Int() int {
	return id.id
}
func ParsePlaybackLogID(id int) (PlaybackLogID, error) {
	if id <= 0 {
		return PlaybackLogID{}, ErrInvalidID
	}

	return PlaybackLogID{id}, nil
}

type PlaybackLog struct {
	id         PlaybackLogID
	trackID    TrackID
	playedAt   time.Time
	amountPaid int
}

func (p PlaybackLog) ID() PlaybackLogID {
	return p.id
}

func (p PlaybackLog) TrackID() TrackID {
	return p.trackID
}

func (p PlaybackLog) PlayedAt() time.Time {
	return p.playedAt
}

func (p PlaybackLog) AmountPaid() int {
	return p.amountPaid
}

var (
	ErrInvalidID         = errors.New("invalid identifier")
	ErrInvalidPaidAmount = errors.New("invalid paid amount")
)

// Note that id is assigned by the database provider.
// Do not use factory created PlaybackLog id for referencing.
func NewPlayBackLog(trackID TrackID, amountPaid int) (*PlaybackLog, error) {
	if amountPaid <= 0 {
		return nil, ErrInvalidPaidAmount
	}

	return &PlaybackLog{
		trackID:    trackID,
		playedAt:   time.Now().UTC(),
		amountPaid: amountPaid,
	}, nil
}

func RebuildPlaybackLog(
	id int,
	trackID int,
	playedAt time.Time,
	amountPaid int,
) (*PlaybackLog, error) {
	playbackLogID, err := ParsePlaybackLogID(id)
	if err != nil {
		return nil, err
	}

	trackIDParsed, err := ParseTrackID(trackID)
	if err != nil {
		return nil, err
	}
	if amountPaid <= 0 {
		return nil, ErrInvalidPaidAmount
	}
	if playedAt.IsZero() {
		return nil, errors.New("invalid playedAt value")
	}

	return &PlaybackLog{
		id:         playbackLogID,
		trackID:    trackIDParsed,
		playedAt:   playedAt,
		amountPaid: amountPaid,
	}, nil
}

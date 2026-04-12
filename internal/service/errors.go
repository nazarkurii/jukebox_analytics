package service

import "errors"

var (
	ErrInvalidTrackID             = errors.New("invalid track id")
	ErrInvalidPaidAmount          = errors.New("invalid paid amount")
	ErrNonExistingTrack           = errors.New("non-existing track")
	ErrInvalidTrackPrice          = errors.New("invalid track price")
	ErrInvalidPopularTracksNumber = errors.New("invalid populars track number")
)

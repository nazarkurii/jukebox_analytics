package model

import (
	"errors"
	"strconv"
)

type TrackID struct {
	id int
}

func (id TrackID) String() string {
	return strconv.Itoa(id.id)
}

func (id TrackID) Int() int {
	return id.id
}
func ParseTrackID(id int) (TrackID, error) {
	if id <= 0 {
		return TrackID{}, ErrInvalidID
	}

	return TrackID{id}, nil
}

type Track struct {
	id     TrackID
	name   string
	artist string
	price  int
}

func (t *Track) ID() TrackID {
	return t.id
}

func (t *Track) Name() string {
	return t.name
}

func (t *Track) Artist() string {
	return t.artist
}

func (t *Track) Price() int {
	return t.price
}

var (
	ErrInvalidTrackPrice = errors.New("invalid track price")
)

func (t *Track) UpdatePrice(newPrice int) error {
	if newPrice <= 0 {
		return ErrInvalidTrackPrice
	}

	t.price = newPrice
	return nil
}

func RebuildTrack(id int, name, artist string, price int) (*Track, error) {
	trackID, err := ParseTrackID(id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, errors.New("invalid track name (empty)")
	}
	if artist == "" {
		return nil, errors.New("invalid track artist (empty)")
	}
	if price <= 0 {
		return nil, ErrInvalidTrackPrice
	}

	return &Track{
		id:     trackID,
		name:   name,
		artist: artist,
		price:  price,
	}, nil
}

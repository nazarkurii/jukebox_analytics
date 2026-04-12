package postresrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/repo"
)

type track struct {
	id     int
	name   string
	artist string
	price  int
}

type Track struct {
	db *pgxpool.Pool
}

func NewTrack(db *pgxpool.Pool) *Track {
	return &Track{db}
}

func (t *Track) GetTrack(ctx context.Context, id model.TrackID) (*model.Track, error) {
	var trackDB track
	err := t.db.QueryRow(
		ctx,
		`SELECT id, name, artist, price FROM tracks WHERE id = $1`,
		id.Int(),
	).Scan(&trackDB.id, &trackDB.name, &trackDB.artist, &trackDB.price)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrNonExistingData
		}
		return nil, err
	}

	track, err := model.RebuildTrack(trackDB.id, trackDB.name, trackDB.artist, trackDB.price)
	if err != nil {
		return nil, fmt.Errorf("failed to rebuild track: %w", err)
	}

	return track, nil
}

func (t *Track) UpdateTrack(ctx context.Context, track *model.Track) error {
	cmd, err := t.db.Exec(
		ctx,
		`UPDATE tracks SET artist = $1, name = $2, price = $3 WHERE id = $4`,
		track.Artist(), track.Name(), track.Price(), track.ID().Int(),
	)

	if cmd.RowsAffected() == 0 {
		return repo.ErrNonExistingData
	}

	return err
}

type Stats struct {
	db *pgxpool.Pool
}

func NewStats(db *pgxpool.Pool) *Stats {
	return &Stats{db}
}

func (s *Stats) GetPopularTracks(ctx context.Context, number int) ([]*model.Track, error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT id, name, artist, price FROM tracks 
		JOIN (SELECT track_id, COUNT(*) times_played FROM playback_logs GROUP BY track_id) popular
			ON popular.track_id = tracks.id
		ORDER BY popular.times_played DESC
		LIMIT $1`,
		number,
	)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Track, 0, number)
	var dbTrack track
	for rows.Next() {
		err = rows.Scan(&dbTrack.id, &dbTrack.name, &dbTrack.artist, &dbTrack.price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		track, err := model.RebuildTrack(dbTrack.id, dbTrack.name, dbTrack.artist, dbTrack.price)
		if err != nil {
			return nil, fmt.Errorf("failed to rebuild track: %w", err)
		}

		result = append(result, track)
	}

	return result, nil
}

func (s *Stats) GetRevenue(ctx context.Context) (result []struct {
	ID           int
	Name, Artist string
	Revenue      int
}, err error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT tracks.id, tracks.name, tracks.artist, revenue.total_paid FROM tracks 
		JOIN (SELECT track_id, SUM(amount_paid) total_paid FROM playback_logs GROUP BY track_id) revenue
			ON revenue.track_id = tracks.id
		ORDER BY revenue.total_paid DESC`,
	)

	if err != nil {
		return nil, err
	}

	var trackName, trackArtist string
	var id, revenue int

	for rows.Next() {
		err := rows.Scan(&id, &trackName, &trackArtist, &revenue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, struct {
			ID           int
			Name, Artist string
			Revenue      int
		}{id, trackName, trackArtist, revenue})
	}

	return result, nil
}

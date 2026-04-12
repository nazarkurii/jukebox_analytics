package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"
)

type statsService interface {
	GetTopTracks(ctx context.Context) ([]*model.Track, error)
	GetRevenue(ctx context.Context) ([]struct {
		ID           int
		Name, Artist string
		Revenue      int
	}, int, error)
}

type Stats struct {
	service statsService
}

func NewStats(service statsService) *Stats {
	return &Stats{service}
}

func (s *Stats) RegisterRoutes(logger *httputil.Logger, mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/stats/top", logger.Wrap(s.getTopTracks))
	mux.HandleFunc("GET /api/v1/stats/revenue", logger.Wrap(s.getRevenue))
}

func (t *Stats) getTopTracks(w http.ResponseWriter, r *http.Request) {
	tracks, err := t.service.GetTopTracks(r.Context())
	if err != nil {
		httputil.WriteInternalServerError(w, fmt.Errorf("failed to get top tracks: %w", err))
		return
	}

	type responseTrack struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Artist string `json:"artist"`
		Price  int    `json:"price"`
	}

	responseTracks := make([]responseTrack, len(tracks))
	for i, track := range tracks {
		responseTracks[i] = responseTrack{
			ID:     track.ID().Int(),
			Name:   track.Name(),
			Artist: track.Artist(),
			Price:  track.Price(),
		}
	}

	httputil.WriteJSON(w, http.StatusOK, responseTracks)
}

func (t *Stats) getRevenue(w http.ResponseWriter, r *http.Request) {
	trackRevenues, totalRevenue, err := t.service.GetRevenue(r.Context())
	if err != nil {
		httputil.WriteInternalServerError(w, fmt.Errorf("failed to get tracks revenue: %w", err))
		return
	}

	type responseTrackRevenue struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Artist  string `json:"artist"`
		Revenue int    `json:"revenue"`
	}

	responseTrackRevenues := make([]responseTrackRevenue, len(trackRevenues))
	for i, trackRevenue := range trackRevenues {
		responseTrackRevenues[i] = responseTrackRevenue(trackRevenue)
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"total_revenue": totalRevenue,
		"by_track":      responseTrackRevenues,
	})
}

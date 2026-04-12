package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/service"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"
)

type playbackService interface {
	RegisterPlayback(ctx context.Context, trackID, amountPaid int) (*model.PlaybackLog, error)
}

type PlaybackLogger struct {
	service playbackService
}

func NewPlaybackLogger(service playbackService) *PlaybackLogger {
	return &PlaybackLogger{service}
}

func (pl *PlaybackLogger) RegisterRoutes(logger *httputil.Logger, mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/logs", logger.Wrap(pl.registerPlayback))
}

type RegisterPlaybackRequest struct {
	TrackID    int `json:"track_id"`
	AmountPaid int `json:"amount_paid"`
}

type ResponsePlaybackLog struct {
	ID         int       `json:"id"`
	TrackID    int       `json:"track_id"`
	PlayedAt   time.Time `json:"played_at"`
	AmountPaid int       `json:"amount_paid"`
}

func (pl *PlaybackLogger) registerPlayback(w http.ResponseWriter, r *http.Request) {
	var request RegisterPlaybackRequest

	if !httputil.UnmarshalJSON(w, r, &request) {
		return
	}

	log, err := pl.service.RegisterPlayback(r.Context(), request.TrackID, request.AmountPaid)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTrackID):
			httputil.WriteErrorJSON(
				w, http.StatusUnprocessableEntity,
				"invalid track id, has to be greater than 0",
				fmt.Errorf("failed register playback log: %w", err),
			)
		case errors.Is(err, service.ErrInvalidPaidAmount):
			httputil.WriteErrorJSON(
				w, http.StatusUnprocessableEntity,
				"invalid amount_paid, has to be greater than 0",
				fmt.Errorf("failed register playback log: %w", err),
			)
		case errors.Is(err, service.ErrNonExistingTrack):
			httputil.WriteNotFoundErrorJSON(w, fmt.Errorf("failed register playback log: %w", err))
		default:
			httputil.WriteInternalServerError(w, fmt.Errorf("failed register playback log: %w", err))
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, ResponsePlaybackLog{
		ID:         log.ID().Int(),
		TrackID:    log.TrackID().Int(),
		PlayedAt:   log.PlayedAt(),
		AmountPaid: log.AmountPaid(),
	})
}

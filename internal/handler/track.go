package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/service"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"
)

type trackService interface {
	UpdatePrice(ctx context.Context, id, newPrice int) (*model.Track, error)
}

type Track struct {
	service trackService
}

func NewTrack(service trackService) *Track {
	return &Track{service}
}

func (t *Track) RegisterRoutes(logger *httputil.Logger, mux *http.ServeMux) {
	mux.HandleFunc("PATCH /api/v1/tracks/{id}/price", logger.Wrap(t.updateTrackPrice))
}

type ResponseTrack struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Price  int    `json:"price"`
}

type UpdateTrackPriceRequest struct {
	NewPrice int `json:"newPrice"`
}

type UpdateTrackPriceResponse struct {
	Track ResponseTrack `json:"track"`
}

func (t *Track) updateTrackPrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httputil.WriteErrorJSON(
			w, http.StatusBadRequest,
			"invalid id format",
			fmt.Errorf("failed to parse track id: %w", err),
		)
		return
	}

	var request UpdateTrackPriceRequest

	if !httputil.UnmarshalJSON(w, r, &request) {
		return
	}

	track, err := t.service.UpdatePrice(r.Context(), id, request.NewPrice)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTrackID):
			httputil.WriteErrorJSON(
				w, http.StatusUnprocessableEntity,
				"invalid track id, has to be greater than 0",
				fmt.Errorf("failed register playback log: %w", err),
			)
		case errors.Is(err, service.ErrInvalidTrackPrice):
			httputil.WriteErrorJSON(
				w, http.StatusBadRequest,
				"invalid track price, has to be greater than 0",
				fmt.Errorf("failed to update track price: %w", err),
			)
		case errors.Is(err, service.ErrNonExistingTrack):
			httputil.WriteNotFoundErrorJSON(w, fmt.Errorf("failed to update track price: %w", err))
		default:
			httputil.WriteInternalServerError(w, fmt.Errorf("failed to update track price: %w", err))
		}

		return
	}

	httputil.WriteJSON(w, http.StatusOK, UpdateTrackPriceResponse{
		ResponseTrack{
			ID:     track.ID().Int(),
			Name:   track.Name(),
			Artist: track.Artist(),
			Price:  track.Price(),
		},
	})
}

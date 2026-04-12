package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/model/testmodel"
	"github.com/nazarkurii/jukebox_analytics/internal/service"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"
)

type playbackServiceMock struct {
	t                             *testing.T
	expectedTrackID, expectedPaid int
	returnLog                     *model.PlaybackLog
	returnErr                     error
}

func (m *playbackServiceMock) init(t *testing.T) {
	m.t = t
}

func (m *playbackServiceMock) RegisterPlayback(ctx context.Context, trackID, amountPaid int) (*model.PlaybackLog, error) {
	var errs error
	if trackID != m.expectedTrackID {
		errs = fmt.Errorf("expected trackID, passed to service.RegisterPlayback, to be %d, got %d", m.expectedTrackID, trackID)
	}

	if amountPaid != m.expectedPaid {
		errs = errors.Join(errs, fmt.Errorf("expected amountPaid, passed to service.RegisterPlayback, to be %d, got %d", m.expectedPaid, amountPaid))
	}

	if errs != nil {
		m.t.Fatal(errs)
	}

	return m.returnLog, m.returnErr
}

func TestPlaybackLogger_registerPlayback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		requestBody  []byte
		expectedErr  *httputil.HandlerError
		expectedCode int
		serviceMock  playbackServiceMock
	}{
		{
			name:        "returns error when json body is invalid",
			requestBody: []byte(`{"track_id": "invalid"}`),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestJsonUnmarshalingClientError(),
				Log:    httputil.TestJsonUnmarshalingLogError[RegisterPlaybackRequest](t, []byte(`{"track_id": "invalid"}`)),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "returns error when track_id is invalid",
			requestBody: []byte(`{"track_id": 0, "amount_paid": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: "invalid track id, has to be greater than 0",
				Log:    "failed register playback log: " + service.ErrInvalidTrackID.Error(),
			},
			serviceMock: playbackServiceMock{
				expectedTrackID: 0,
				expectedPaid:    100,
				returnErr:       service.ErrInvalidTrackID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "returns error when amount_paid is invalid",
			requestBody: []byte(`{"track_id": 1, "amount_paid": -10}`),
			expectedErr: &httputil.HandlerError{
				Client: "invalid amount_paid, has to be greater than 0",
				Log:    "failed register playback log: " + service.ErrInvalidPaidAmount.Error(),
			},
			serviceMock: playbackServiceMock{
				expectedTrackID: 1,
				expectedPaid:    -10,
				returnErr:       service.ErrInvalidPaidAmount,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "returns error when track does not exist",
			requestBody: []byte(`{"track_id": 999, "amount_paid": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestNonExistingRessourceErr(),
				Log:    "failed register playback log: " + service.ErrNonExistingTrack.Error(),
			},
			serviceMock: playbackServiceMock{
				expectedTrackID: 999,
				expectedPaid:    100,
				returnErr:       service.ErrNonExistingTrack,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "returns error when service.RegisterPlayback fails unexpectedly",
			requestBody: []byte(`{"track_id": 1, "amount_paid": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestInternalServerClientErr(),
				Log:    "failed register playback log: internal error",
			},
			serviceMock: playbackServiceMock{
				expectedTrackID: 1,
				expectedPaid:    100,
				returnErr:       errors.New("internal error"),
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "returns playback log when succeeds",
			requestBody: []byte(`{"track_id": 1, "amount_paid": 100}`),
			serviceMock: playbackServiceMock{
				expectedTrackID: 1,
				expectedPaid:    100,
				returnLog:       testmodel.TestPlaybackLog(t),
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.serviceMock.init(t)

			req := httptest.NewRequest(
				"POST",
				"/api/v1/logs",
				bytes.NewReader(tc.requestBody),
			)
			rec := httptest.NewRecorder()

			NewPlaybackLogger(&tc.serviceMock).registerPlayback(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("expected response code to be %d, got %d", tc.expectedCode, rec.Code)
			}

			if tc.expectedErr != nil {
				httputil.TestErrorCheck(t, rec.Body.Bytes(), *tc.expectedErr)
				return
			}

			var response ResponsePlaybackLog
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("failed to unmarshal recorder response: %s", err)
			}

			expectedResponse := ResponsePlaybackLog{
				ID:         tc.serviceMock.returnLog.ID().Int(),
				TrackID:    tc.serviceMock.returnLog.TrackID().Int(),
				PlayedAt:   tc.serviceMock.returnLog.PlayedAt(),
				AmountPaid: tc.serviceMock.returnLog.AmountPaid(),
			}

			// Use Equal for time comparison to avoid monotonic clock issues
			if !response.PlayedAt.Equal(expectedResponse.PlayedAt) {
				t.Errorf("expected PlayedAt %v, got %v", expectedResponse.PlayedAt, response.PlayedAt)
			}

			// Clear time fields for simpler struct comparison
			response.PlayedAt = expectedResponse.PlayedAt
			if response != expectedResponse {
				t.Fatalf("expected response %+v, got %+v", expectedResponse, response)
			}
		})
	}
}

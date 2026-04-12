package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/nazarkurii/jukebox_analytics/internal/model"
	"github.com/nazarkurii/jukebox_analytics/internal/model/testmodel"
	"github.com/nazarkurii/jukebox_analytics/internal/service"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"
)

type trackServiceMock struct {
	t                            *testing.T
	expectedID, expectedNewPrice int
	returnTrack                  *model.Track
	returnErr                    error
}

func (m *trackServiceMock) init(t *testing.T) {
	m.t = t
}

func (m *trackServiceMock) UpdatePrice(ctx context.Context, id, newPrice int) (*model.Track, error) {
	var errs error
	if id != m.expectedID {
		errs = fmt.Errorf("expected id, passed to service.UpdatePrice, to be %q, got %q", m.expectedID, id)
	}

	if newPrice != m.expectedNewPrice {
		errs = errors.Join(errs, fmt.Errorf("expected newPrice, passed to service.UpdatePrice, to be %q, got %q", m.expectedNewPrice, newPrice))
	}

	if errs != nil {
		m.t.Fatal(errs)
	}

	return m.returnTrack, m.returnErr
}

func TestTrack_updateTrackPrice(t *testing.T) {
	t.Parallel()

	getAtoiErr := func(from string) string {
		_, err := strconv.Atoi(from)
		if err == nil {
			t.Fatal("failed to get atoi error")
		}
		return err.Error()
	}

	tests := []struct {
		name         string
		id           string
		requestBody  []byte
		expectedErr  *httputil.HandlerError
		expectedCode int
		serviceMoc   trackServiceMock
	}{
		{
			name: "returns error when id is invalid",
			id:   "non-number",
			expectedErr: &httputil.HandlerError{
				Client: "invalid id format",
				Log:    "failed to parse track id: " + getAtoiErr("non-number"),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "returns error when id is invalid",
			id:   "1.1",
			expectedErr: &httputil.HandlerError{
				Client: "invalid id format",
				Log:    "failed to parse track id: " + getAtoiErr("1.1"),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "returns error when json body is invalid",
			id:          "1",
			requestBody: []byte("{newPrice: invalid}"),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestJsonUnmarshalingClientError(),
				Log:    httputil.TestJsonUnmarshalingLogError[UpdateTrackPriceRequest](t, []byte("{newPrice: invalid}")),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "returns error when track id is invalid at service level",
			id:          "0",
			requestBody: []byte(`{"newPrice": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: "invalid track id, has to be greater than 0",
				Log:    "failed register playback log: " + service.ErrInvalidTrackID.Error(),
			},
			serviceMoc: trackServiceMock{
				expectedID:       0,
				expectedNewPrice: 100,
				returnErr:        service.ErrInvalidTrackID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "returns error when id is not valid",
			id:          "-1",
			requestBody: []byte(`{"newPrice": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: "invalid track price, has to be greater than 0",
				Log:    "failed to update track price: " + service.ErrInvalidTrackPrice.Error(),
			},
			serviceMoc: trackServiceMock{
				expectedID:       -1,
				expectedNewPrice: 100,
				returnErr:        service.ErrInvalidTrackPrice,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "returns error when track does not exist",
			id:          "1",
			requestBody: []byte(`{"newPrice": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestNonExistingRessourceErr(),
				Log:    "failed to update track price: " + service.ErrNonExistingTrack.Error(),
			},
			serviceMoc: trackServiceMock{
				expectedID:       1,
				expectedNewPrice: 100,
				returnErr:        service.ErrNonExistingTrack,
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "returns error when service.UpdatePrice fails unexpectedly",
			id:          "1",
			requestBody: []byte(`{"newPrice": 100}`),
			expectedErr: &httputil.HandlerError{
				Client: httputil.TestInternalServerClientErr(),
				Log:    "failed to update track price: internal error",
			},
			serviceMoc: trackServiceMock{
				expectedID:       1,
				expectedNewPrice: 100,
				returnErr:        errors.New("internal error"),
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:        "returns track with updated price when succeeds",
			id:          "1",
			requestBody: []byte(`{"newPrice": 100}`),
			serviceMoc: trackServiceMock{
				expectedID:       1,
				expectedNewPrice: 100,
				returnTrack:      testmodel.TestTrack(t),
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.serviceMoc.init(t)

			req := httptest.NewRequest(
				"PATCH",
				"/api/v1/tracks/{id}/price",
				bytes.NewReader(tc.requestBody),
			)

			req.SetPathValue("id", tc.id)
			rec := httptest.NewRecorder()

			NewTrack(&tc.serviceMoc).updateTrackPrice(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("expected response code to be %d, got %d", tc.expectedCode, rec.Code)
			}

			if tc.expectedErr != nil {
				httputil.TestErrorCheck(t, rec.Body.Bytes(), *tc.expectedErr)
				return
			}

			var response UpdateTrackPriceResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("failed to unmarshal recorder response: %s", err)
			}

			expectedResponse := UpdateTrackPriceResponse{
				ResponseTrack{
					ID:     tc.serviceMoc.returnTrack.ID().Int(),
					Name:   tc.serviceMoc.returnTrack.Name(),
					Artist: tc.serviceMoc.returnTrack.Artist(),
					Price:  tc.serviceMoc.returnTrack.Price(),
				},
			}

			if response != expectedResponse {
				t.Fatalf("expected resposne to be %q, got %q", expectedResponse, response)
			}

		})
	}
}

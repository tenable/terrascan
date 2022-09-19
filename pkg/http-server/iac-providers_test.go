package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIacProviders(t *testing.T) {

	handler := NewAPIHandler()

	t.Run("test providers api status", func(t *testing.T) {
		var (
			req, _ = http.NewRequest(http.MethodGet, "/v1/providers", nil)
			resp   = httptest.NewRecorder()
			want   = http.StatusOK
		)
		handler.iacProviders(resp, req)
		got := resp.Result().StatusCode

		if got != want {
			t.Errorf("incorrect providers status code, got: '%v', want: '%v'", got, want)
		}
	})

	t.Run("test providers api response structure", func(t *testing.T) {
		var (
			req, _ = http.NewRequest(http.MethodGet, "/v1/providers", nil)
			resp   = httptest.NewRecorder()
		)
		handler.iacProviders(resp, req)

		var data []IacProvider
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			t.Errorf("error parsing response body, error: '%v'", err.Error())
		}
	})
}

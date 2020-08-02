package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {

	handler := NewAPIHandler()

	t.Run("test health api", func(t *testing.T) {
		var (
			req, _ = http.NewRequest(http.MethodGet, "/players/Pepper", nil)
			resp   = httptest.NewRecorder()
			want   = http.StatusOK
		)
		handler.Health(resp, req)
		got := resp.Result().StatusCode

		if got != want {
			t.Errorf("incorrect health status code, got: '%v', want: '%v'", got, want)
		}
	})
}

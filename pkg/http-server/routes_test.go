package httpserver

import (
	"testing"
)

func TestRoutes(t *testing.T) {
	t.Run("health route check", func(t *testing.T) {
		var (
			server = NewAPIServer()
			got    = server.Routes()
			passed = false
		)

		for _, r := range got {
			if r.path == "/health" && r.verb == "GET" {
				passed = true
				break
			}
		}
		if !passed {
			t.Errorf("failed to find /health in routes")
		}
	})
}

func TestPath(t *testing.T) {

	table := []struct {
		name    string
		route   string
		version string
		want    string
	}{
		{
			name:    "route someroute version v1",
			route:   "/someroute",
			version: "v1",
			want:    "/v1/someroute",
		},
	}

	for _, tt := range table {
		got := versionedPath(tt.route)
		if got != tt.want {
			t.Errorf("got: '%v', want: '%v'", got, tt.want)
		}
	}
}

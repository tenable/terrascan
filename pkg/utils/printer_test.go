package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {

	table := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "empty JSON",
			input: make(map[string]interface{}),
			want:  "{}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := &bytes.Buffer{}
			PrintJSON(tt.input, got)
			if strings.TrimSpace(got.String()) != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

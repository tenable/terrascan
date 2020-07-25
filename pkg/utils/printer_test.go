package utils

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	validJSONFile = "./testdata/valid.json"
)

var (
	validJSON      []byte
	validJSONInput = map[string]int{"apple": 5, "lettuce": 7}
)

func init() {
	validJSON, _ = ioutil.ReadFile(validJSONFile)

}

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
		{
			name:  "valid JSON",
			input: validJSONInput,
			want:  string(validJSON),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := &bytes.Buffer{}
			PrintJSON(tt.input, got)
			if strings.TrimSpace(got.String()) != strings.TrimSpace(tt.want) {
				t.Errorf("got:\n'%v'\n, want:\n'%v'\n", got, tt.want)
			}
		})
	}
}

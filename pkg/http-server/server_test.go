package httpserver

import (
	"reflect"
	"testing"
)

func TestNewAPIServer(t *testing.T) {
	t.Run("new API gateway", func(t *testing.T) {
		var (
			want = APIServer{}
			got  = NewAPIServer()
		)
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("got: '%v', want: '%v'", *got, want)
		}
	})
}

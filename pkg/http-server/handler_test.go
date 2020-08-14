package httpserver

import (
	"reflect"
	"testing"
)

func TestNewAPIHandler(t *testing.T) {
	t.Run("new API gateway", func(t *testing.T) {
		var (
			want = APIHandler{}
			got  = NewAPIHandler()
		)
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("got: '%v', want: '%v'", *got, want)
		}
	})
}

package httpserver

import (
	"reflect"
	"testing"
)

func TestNewAPIGateway(t *testing.T) {
	t.Run("new API gateway", func(t *testing.T) {
		var (
			want = APIGateway{}
			got  = NewAPIGateway()
		)
		if !reflect.DeepEqual(*got, want) {
			t.Errorf("got: '%v', want: '%v'", *got, want)
		}
	})
}

package termcolor

import (
	"strings"
	"testing"
)

func TestColoring(t *testing.T) {
	s := "FORCED FAILURE"
	s = Colorize("Fg#fff|Bg#f00|Bold", s)
	if !strings.HasPrefix(s, ColorPrefix) {
		t.Errorf("expected console escape code %v, got %v", []byte(ColorPrefix), []byte(s[:len(ColorPrefix)]))
	}
}

package cli

import (
	"github.com/accurics/terrascan/pkg/termcolor"
	"io"
	"os"
)

func NewOutputWriter(useColors bool) io.Writer {

	// Color codes will corrupt output, so suppress if not on terminal
	if useColors == true {
		return termcolor.NewColorizedWriter(os.Stdout)
	}
	return os.Stdout
}

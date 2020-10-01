package cli

import (
	"github.com/accurics/terrascan/pkg/termcolor"
	"io"
	"os"
)

// NewOutputWriter gets a new io.Writer based on os.Stdout.
// If param useColors=true, the writer will colorize the output
func NewOutputWriter(useColors bool) io.Writer {

	// Color codes will corrupt output, so suppress if not on terminal
	if useColors {
		return termcolor.NewColorizedWriter(os.Stdout)
	}
	return os.Stdout
}

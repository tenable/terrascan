package termcolor

import (
	"fmt"
	"io"
)

// ColorizedWriter provides an io.Writer which colorizes output
type ColorizedWriter struct {
	writer io.Writer
}

// NewColorizedWriter returns a new io.Writer which implements coloring
func NewColorizedWriter(w io.Writer) ColorizedWriter {
	return ColorizedWriter{w}
}

func (me ColorizedWriter) Write(p []byte) (n int, err error) {
	/* Before output is written, perform color substitutions */
	for ptn, style := range GetColorPatterns() {
		p = ptn.ReplaceAllFunc(p, func(matched []byte) []byte {
			groups := ptn.FindStringSubmatch(string(matched))
			return []byte(fmt.Sprintf("%s%s%s%s%s",
				groups[1],
				Colorize(style.KeyStyle, string(groups[2])),
				groups[3],
				Colorize(style.ValueStyle, string(groups[4])),
				groups[5]))
		})
	}

	return me.writer.Write(p)
}

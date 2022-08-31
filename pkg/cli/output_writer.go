package cli

import (
	"io"
	"os"
	"path/filepath"

	"github.com/tenable/terrascan/pkg/termcolor"
	"go.uber.org/zap"
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

// NewFileWriter gets a new io.Writer based on file output and closing function.
// It returns `nil nil` if the value of dir is empty or if the file can't be opened
func NewFileWriter(dir string, outputType string) (io.Writer, func() error) {

	// if no directory resolved
	if dir == "" {
		return nil, nil
	}

	fileName := "scan-result.txt"

	// decide the file extension/type
	switch outputType {
	case "json":
		fileName = "scan-result.json"
	case "yaml":
		fileName = "scan-result.yaml"
	case "xml", "junit-xml":
		fileName = "scan-result.xml"
	case "sarif", "github-sarif":
		fileName = "scan-result.sarif"
	}

	filePath := filepath.Join(dir, fileName)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		zap.S().Warn("unable to open file: %s, error:%w", filePath, err)
		return nil, nil
	}

	return f, f.Close
}

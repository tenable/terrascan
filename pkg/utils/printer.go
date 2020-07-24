package utils

import (
	"encoding/json"
	"io"

	"go.uber.org/zap"
)

// PrintJSON prints data in JSON format
func PrintJSON(data interface{}, writer io.Writer) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		zap.S().Errorf("failed to create JSON. error: '%v'", err)
		return
	}
	writer.Write(j)
	writer.Write([]byte{'\n'})
}

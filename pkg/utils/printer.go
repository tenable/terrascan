package utils

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

// PrintJSON prints data in JSON format
func PrintJSON(data interface{}) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		zap.S().Errorf("failed to create JSON. error: '%v'", err)
		return
	}
	os.Stdout.Write(j)
	os.Stdout.Write([]byte{'\n'})
}

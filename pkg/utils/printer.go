package utils

import (
	"encoding/json"
	"log"
	"os"
)

// PrintJSON prints data in JSON format
func PrintJSON(data interface{}) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("failed to create JSON. error: '%v'", err)
		return
	}
	os.Stdout.Write(j)
	os.Stdout.Write([]byte{'\n'})
}

package utils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
)

const (
	// JSONDoc type for json files
	JSONDoc = "json"
)

// LoadJSON loads a JSON file into an IacDocument struct
func LoadJSON(filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 1)

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return iacDocumentList, err
	}

	// First pass determines line number data
	currentLineNumber := 1
	{ // Limit the scope for Close()
		var file *os.File
		file, err = os.Open(filePath)
		if err != nil {
			return iacDocumentList, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			currentLineNumber++
		}

		if err = scanner.Err(); err != nil {
			return iacDocumentList, err
		}
	}

	// Second pass extracts data
	var doc IacDocument
	dataMap := make(map[string]interface{})
	err = json.Unmarshal(fileBytes, &dataMap)
	if err != nil {
		zap.S().Warn("unable to unmarshal json file", zap.String("file", filePath))
		return iacDocumentList, err
	}

	doc.Data, err = json.Marshal(dataMap)
	if err != nil {
		zap.S().Warn("unable to marshal json file", zap.String("file", filePath))
	}

	doc.StartLine = 1
	doc.FilePath = filePath
	doc.EndLine = currentLineNumber
	iacDocumentList[0] = &doc

	return iacDocumentList, nil
}

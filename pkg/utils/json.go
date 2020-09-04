package utils

import (
	"bufio"
	"io/ioutil"
	"os"
)

const (
	// JSONDoc type for json files
	JSONDoc = "json"
)

// LoadJSON loads a JSON file into an IacDocument struct
func LoadJSON(filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 0, 1)

	// First pass determines line number data
	currentLineNumber := 1
	{ // Limit the scope for Close()
		file, err := os.Open(filePath)
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

	// Second pass extracts raw data
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return iacDocumentList, err
	}

	iacDocumentList = append(iacDocumentList, &IacDocument{
		Type:      JSONDoc,
		StartLine: 1,
		FilePath:  filePath,
		EndLine:   currentLineNumber,
		Data:      fileBytes,
	})

	return iacDocumentList, nil
}

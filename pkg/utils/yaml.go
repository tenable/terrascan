package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"
)

const (
	// YAMLDoc type for yaml files
	YAMLDoc = "yaml"
)

var (
	errHighDocumentCount = fmt.Errorf("document count was higher than expected count")
)

// LoadYAML loads a YAML file. Can return one or more IaC Documents.
func LoadYAML(filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 0)

	// First pass determines line number data
	{ // Limit the scope for Close()
		file, err := os.Open(filePath)
		if err != nil {
			return iacDocumentList, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		startLineNumber := 1
		currentLineNumber := 1
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "---") {
				// We've found the end-of-directives marker, so record results for the current document
				iacDocumentList = append(iacDocumentList, &IacDocument{
					Type:      YAMLDoc,
					StartLine: startLineNumber,
					EndLine:   currentLineNumber,
					FilePath:  filePath,
				})
				startLineNumber = currentLineNumber + 1
			}
			currentLineNumber++
		}

		// Add the very last entry
		iacDocumentList = append(iacDocumentList, &IacDocument{
			Type:      YAMLDoc,
			StartLine: startLineNumber,
			EndLine:   currentLineNumber,
			FilePath:  filePath,
		})

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Second pass extracts all YAML documents and saves it in the document struct
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return iacDocumentList, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(fileBytes))
	i := 0
	for {
		// each iteration extracts and marshals one yaml document
		var value interface{}
		err = dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return iacDocumentList, err
		}
		if i > (len(iacDocumentList) - 1) {
			return iacDocumentList, errHighDocumentCount
		}

		var documentBytes []byte
		documentBytes, err = yaml.Marshal(value)
		if err != nil {
			return iacDocumentList, err
		}
		iacDocumentList[i].Data = documentBytes
		i++
	}

	return iacDocumentList, nil
}

// YAMLtoJSON converts YAML byte data to JSON bytes
func YAMLtoJSON(data []byte) (*map[string]interface{}, error) {
	// fetch the YAML data into an interface type
	var dataMap interface{}
	err := yaml.Unmarshal(data, &dataMap)
	if err != nil {
		zap.S().Warn("unable to unmarshal yaml data")
		return nil, err
	}

	// convert map[interface]interface to map[string]interface throughout the YAML data
	var dataBytes []byte
	dataMap = InterfaceToMapStringInterface(dataMap)

	// marshal to json to produce the json bytes
	if dataBytes, err = json.Marshal(dataMap); err != nil {
		zap.S().Warn("unable to marshal json during conversion")
		return nil, err
	}

	// convert back to map[string]interface with the json data
	configData := make(map[string]interface{})
	if err = json.Unmarshal(dataBytes, &configData); err != nil {
		zap.S().Warn("unable to unmarshal json during conversion")
		return nil, err
	}

	return &configData, nil
}

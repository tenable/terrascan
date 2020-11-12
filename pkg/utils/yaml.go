package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
// Besides reading in file data, its main purpose is to determine and store line number and filename metadata
func LoadYAML(filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return iacDocumentList, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return iacDocumentList, err
	}

	return getIacDocumentList(bufio.NewScanner(file), fileBytes, filePath)
}

// LoadYAMLString loads a YAML String. Can return one or more IaC Documents.
// Besides reading in file data, its main purpose is to determine and store line number and filename metadata
func LoadYAMLString(data, absFilePath string) ([]*IacDocument, error) {
	return getIacDocumentList(bufio.NewScanner(strings.NewReader(data)), []byte(data), absFilePath)
}

// ReadYamlFile reads a yaml file and load content in a map[string]interface{} type
func ReadYamlFile(path string) (map[string]interface{}, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	err = yaml.Unmarshal(dat, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// getIacDocumentList provides one or more IaC Documents.
// Besides reading in file data, its main purpose is to determine and store line number and filename metadata
func getIacDocumentList(scanner *bufio.Scanner, bytearray []byte, filePath string) ([]*IacDocument, error) {
	iacDocumentList := make([]*IacDocument, 0)

	// First pass determines line number data
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
		FilePath:  "",
	})

	if err := scanner.Err(); err != nil {
		return iacDocumentList, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(bytearray))

	i := 0
	for {
		// each iteration extracts and marshals one yaml document
		var value interface{}
		err := dec.Decode(&value)
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

package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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

// AreEqualJSON validate if two json strings are equal
func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	errmsg := "error json unmarshalling string: %s. error: %v"

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf(errmsg, s1, err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf(errmsg, s2, err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// AreEqualJSONBytes validate if two json byte arrays are equal
func AreEqualJSONBytes(b1, b2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	errmsg := "error json unmarshalling bytes: %s. error: %v"

	var err error
	err = json.Unmarshal(b1, &o1)
	if err != nil {
		return false, fmt.Errorf(errmsg, b1, err.Error())
	}
	err = json.Unmarshal(b2, &o2)
	if err != nil {
		return false, fmt.Errorf(errmsg, b2, err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

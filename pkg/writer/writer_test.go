package writer

import (
	"bytes"
	"testing"
)

func TestWriteWithNonRegisteredWriter(t *testing.T) {
	err := Write("test", nil, &bytes.Buffer{})
	if err != errNotSupported {
		t.Errorf("Expected error = %v but got %v", errNotSupported, err)
	}
}

func TestWriteWithRegisteredWriter(t *testing.T) {
	err := Write("json", nil, &bytes.Buffer{})
	if err != nil {
		t.Errorf("Unexpected error for json writer = %v", err)
	}
}

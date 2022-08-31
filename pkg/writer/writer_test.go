package writer

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteWithNonRegisteredWriter(t *testing.T) {
	var bf bytes.Buffer
	w := []io.Writer{&bf}
	err := Write("test", nil, w)
	if err != errNotSupported {
		t.Errorf("Expected error = %v but got %v", errNotSupported, err)
	}
}

func TestWriteWithRegisteredWriter(t *testing.T) {
	var bf bytes.Buffer
	w := []io.Writer{&bf}
	err := Write("json", nil, w)
	if err != nil {
		t.Errorf("Unexpected error for json writer = %v", err)
	}
}

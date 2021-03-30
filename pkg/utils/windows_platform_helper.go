package utils

import (
	"bytes"
	"runtime"
	"strings"
)

// IsWindowsPlatform checks if os is windows
func IsWindowsPlatform() bool {
	return runtime.GOOS == "windows"
}

// ReplaceWinNewLineBytes replaces windows new lines with unix new lines in a byte slice
func ReplaceWinNewLineBytes(input []byte) []byte {
	return bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n"))
}

// ReplaceWinNewLineString replaces windows new lines with unix new lines in a string
func ReplaceWinNewLineString(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}

// ReplaceCarriageReturnBytes replaces windows new lines characters in a string
func ReplaceCarriageReturnBytes(input []byte) []byte {
	return bytes.ReplaceAll(input, []byte("\\r\\n"), []byte("\\n"))
}

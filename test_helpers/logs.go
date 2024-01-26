package test_helpers

import (
	"bytes"
	"log"
)

// CaptureOutput redirects the log output to a buffer and returns a function to restore the original state and the buffer.
func CaptureOutput() (*bytes.Buffer, func()) {
	originalOutput := log.Writer() // Store the original output
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	return buf, func() {
		log.SetOutput(originalOutput) // Restore the original output
	}
}

package codec_test

import (
	"errors"
)

type errorReader struct{}

func (reader errorReader) Read([]byte) (int, error) {
	return 10, errors.New("error")
}

type errorWriter struct{}

func (writer errorWriter) Write([]byte) (int, error) {
	return 10, errors.New("error")
}

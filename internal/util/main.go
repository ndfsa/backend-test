package util

import (
	"encoding/json"
	"errors"
	"io"
)

func DecodeJson[T any](body io.ReadCloser, cl *T) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(cl); err != nil {
		return errors.New("invalid parameters")
	}

	return nil
}


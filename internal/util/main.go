package util

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Send[T any](w *http.ResponseWriter, data T) error {
	(*w).Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(*w).Encode(data); err != nil {
		return err
	}
	return nil
}

func Receive[T any](body io.ReadCloser, cl *T) error {
	if err := json.NewDecoder(body).Decode(cl); err != nil {
		return err
	}
	return nil
}

func SendError(w *http.ResponseWriter, status int, message string) {
	(*w).WriteHeader(status)
	log.Println(message)
}

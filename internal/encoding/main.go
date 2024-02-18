package encoding

import (
	"encoding/json"
	"net/http"
)

func Send[T any](w http.ResponseWriter, data T) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func Receive[T any](r *http.Request, cl *T) error {
	if err := json.NewDecoder(r.Body).Decode(cl); err != nil {
		return err
	}
	return nil
}

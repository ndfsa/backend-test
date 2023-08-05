package internal

import (
	"log"
	"net/http"
)

func handleError(w *http.ResponseWriter, err error, message string) {
	log.Println(err.Error())
	http.Error(*w, message, )

}

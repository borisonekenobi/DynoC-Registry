package handlers

import (
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

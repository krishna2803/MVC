package controller

import (
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong 🏓")
	w.WriteHeader(http.StatusOK)
}

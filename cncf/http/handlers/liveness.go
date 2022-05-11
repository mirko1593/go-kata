package handlers

import "net/http"

// Liveness ...
func Liveness(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

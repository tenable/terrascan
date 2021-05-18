package httputils

import "net/http"

// NotFound writes 404 to the write object
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// NotAllowed writes 405 to the writer object
func NotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

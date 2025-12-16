package utils

import (
	"fmt"
	"net/http"
)

func ReponseError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf(`{"error": %q}`, err), code)
}

func ReponseStatus(w http.ResponseWriter, body []byte, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.Write(body)
}

package utils

import "net/http"

func WriteInternalServerError(w *http.ResponseWriter, e error) {
	(*w).WriteHeader(http.StatusInternalServerError)
	(*w).Write([]byte(e.Error()))
}

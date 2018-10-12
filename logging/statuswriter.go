package logging

import (
	"net/http"
)

type StatusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w StatusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w StatusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (w StatusWriter) Status() int {
	if w.status == 0 {
		w.status = 200
	}
	return w.status
}

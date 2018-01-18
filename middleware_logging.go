package main

import (
	"net/http"
	"time"
)

// ResponseWriter wrapper that keeps track of the status code we sent
type responseWriterWithStatus struct {
	http.ResponseWriter
	code int
}

func (w *responseWriterWithStatus) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *responseWriterWithStatus) Status() int { return w.code }

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wWithStatus := &responseWriterWithStatus{w, 0}
		start := time.Now()
		h.ServeHTTP(wWithStatus, r)
		log.WithFields(log.Fields{
			"method": r.Method,
			"code":   wWithStatus.Status(),
			"ms":     time.Now().Unix() - start.Unix(),
		}).Info(r.URL.Path)
	})
}

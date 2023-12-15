package internalhttp

import (
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	time.Time
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		ResponseWriter: w,
		Time:           time.Now(),
		statusCode:     http.StatusOK,
	}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

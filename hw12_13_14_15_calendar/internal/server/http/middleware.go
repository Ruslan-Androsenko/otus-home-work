package internalhttp

import (
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		latency := time.Since(lrw.Time)
		userAgent := r.Header.Get("User-Agent")

		logg.Infof("%s %s %s %s %d %s %v %s",
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto,
			statusCode, http.StatusText(statusCode), latency, userAgent)
	})
}

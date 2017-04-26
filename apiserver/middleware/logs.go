package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logs is a middleware function that logs every request to the server and it's response time
func Logs(logger *log.Logger) Adapter {
	//return an Adapter function that...
	return func(handler http.Handler) http.Handler {
		//returns an http.Handler that...
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := NewLoggingResponseWriter(w)
			handler.ServeHTTP(lrw, r)
			end := time.Since(start)
			statusCode := lrw.statusCode
			logger.Printf("%v %v [%v] %v %v %d %s", start, r.RemoteAddr, r.Method, r.URL, end, statusCode, http.StatusText(statusCode))
		})
	}
}

// LoggingResponseWriter is a wrapper struct for a responce writer
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter is a wrapper that accesses the underlying response code from a response Writer
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader is a method that writes a response code to a ResponseWriter
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

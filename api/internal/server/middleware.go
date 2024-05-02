package server

import (
	"net/http"

	"github.com/furrygem/ipgem/api/internal/logger"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

type LoggingMiddleware struct{}

func (lm *LoggingMiddleware) Next(next http.Handler) http.Handler {
	l := logger.GetLogger()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		lrw := LoggingResponseWriter{
			ResponseWriter: rw,
			statusCode:     200,
		}
		next.ServeHTTP(&lrw, r)
		l.Infof("%s %s %d", r.Method, r.URL, lrw.statusCode)
	})
}

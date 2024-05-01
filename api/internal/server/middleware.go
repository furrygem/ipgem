package server

import (
	"net/http"

	"github.com/furrygem/ipgem/api/internal/logger"
)

type Middleware struct{}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	l := logger.GetLogger()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		l.Infof("%s %s", r.Method, r.URL)
		next.ServeHTTP(rw, r)
	})
}

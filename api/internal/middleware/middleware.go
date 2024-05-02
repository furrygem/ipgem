package middleware

import "net/http"

type Middleware interface {
	Next(http.Handler) http.Handler
}

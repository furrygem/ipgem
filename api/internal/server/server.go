package server

import (
	"net/http"

	"github.com/furrygem/ipgem/api/internal/middleware"
)

type Server struct {
	listenAddress  string
	router         *http.ServeMux
	middlewareList []middleware.Middleware
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	s.registerRoutes()
	serv := s.addMiddleware([]middleware.Middleware{&LoggingMiddleware{}})
	return http.ListenAndServe(s.listenAddress, serv)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/hello", helloHandler)
}

func (s *Server) addMiddleware(middlewareList []middleware.Middleware) http.Handler {
	var serv http.Handler
	serv = s
	for _, m := range middlewareList {
		serv = m.Next(serv)
	}
	return serv
}

func NewServer(addr string) *Server {
	router := http.NewServeMux()
	return &Server{
		router:        router,
		listenAddress: addr,
	}
}

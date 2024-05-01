package server

import "net/http"

type Server struct {
	listenAddress string
	router        *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	s.registerRoutes()
	middleware := Middleware{}
	return http.ListenAndServe(s.listenAddress, middleware.Logging(s))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/hello", helloHandler)
}

func NewServer(addr string) *Server {
	router := http.NewServeMux()
	return &Server{
		router:        router,
		listenAddress: addr,
	}
}

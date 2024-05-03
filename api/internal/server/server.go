package server

import (
	"encoding/json"
	"net/http"

	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/middleware"
	"github.com/furrygem/ipgem/api/internal/repository"
	"github.com/furrygem/ipgem/api/internal/service"
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
	repo, err := repository.NewSQLiteRepository("db.sqlite3")
	logger := logger.GetLogger()
	if err != nil {
		logger.Fatal(err)
	}
	s := service.NewService(repo)
	defer s.CloseConn()
	res, err := s.ListRecords()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}
	jsoned, err := json.Marshal(res)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(jsoned)
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

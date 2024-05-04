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

type Handler struct {
	service *service.DNSCrud
}

func NewHandler(repo repository.Repository) *Handler {
	service := service.NewService(repo)
	return &Handler{
		service: service,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	l := logger.GetLogger()
	repo, err := repository.NewSQLiteRepository("db.sqlite3")
	if err != nil {
		l.Fatal(err)
	}
	err = repo.Open()
	if err != nil {
		l.Fatal(err)
	}
	defer repo.Close()
	handler := NewHandler(repo)
	s.registerRoutes(handler)
	serv := s.addMiddleware([]middleware.Middleware{&LoggingMiddleware{}})
	return http.ListenAndServe(s.listenAddress, serv)
}

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()
	res, err := h.service.ListRecords()
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

func (h *Handler) retrieveHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.GetLogger()
	queryID := r.PathValue("id")
	if queryID == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	record, err := h.service.RetrieveRecord(queryID)
	if err != nil {
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	jsoned, err := json.Marshal(record)
	if err != nil {
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	rw.Write(jsoned)

}

func (s *Server) registerRoutes(h *Handler) {
	s.router.HandleFunc("/records", h.listHandler)
	s.router.HandleFunc("/records/{id}", h.retrieveHandler)
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

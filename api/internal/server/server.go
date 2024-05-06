package server

import (
	"encoding/json"
	"net/http"

	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/middleware"
	"github.com/furrygem/ipgem/api/internal/models"
	"github.com/furrygem/ipgem/api/internal/repository"
	"github.com/furrygem/ipgem/api/internal/service"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

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
	validate = validator.New()
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
	// FIXME: Validate its a valid uuid
	if queryID == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	record, err := h.service.RetrieveRecord(queryID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			l.Warn(err)
			rw.WriteHeader(404)
			rw.Write([]byte("Not found"))
			return
		}
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

func (h *Handler) updateHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.GetLogger()
	// FIXME: Validate its a valid uuid
	queryID := r.PathValue("id")
	if queryID == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	recordDTO := models.RecordDTO{}
	err := json.NewDecoder(r.Body).Decode(&recordDTO)
	if err != nil {
		l.Warn(err)
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	err = validate.Struct(recordDTO)
	if err != nil {
		l.Warn(err)
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	updatedRecord, err := h.service.UpdateRecord(queryID, &recordDTO)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			l.Warn(err)
			rw.WriteHeader(404)
			rw.Write([]byte("Not found"))
			return
		}
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	jsoned, err := json.Marshal(updatedRecord)
	if err != nil {
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	rw.Write(jsoned)
}

func (h *Handler) createHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.GetLogger()
	recordDTO := &models.RecordDTO{}
	err := json.NewDecoder(r.Body).Decode(recordDTO)
	if err != nil {
		l.Warn(err)
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	validate.Struct(recordDTO)
	inserted, err := h.service.AddRecord(recordDTO)
	if err != nil {
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	jsoned, err := json.Marshal(inserted)
	if err != nil {
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	rw.Write(jsoned)
}

func (h *Handler) deleteHandler(rw http.ResponseWriter, r *http.Request) {
	l := logger.GetLogger()
	queryID := r.PathValue("id")
	if queryID == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	err := h.service.DeleteRecord(queryID)
	if err != nil {
		if err.Error() == "No rows affected by delete" {
			l.Warn(err)
			rw.WriteHeader(404)
			rw.Write([]byte("Not found"))
			return
		}
		l.Error(err)
		rw.WriteHeader(500)
		rw.Write([]byte("Internal server error"))
		return
	}
	rw.WriteHeader(http.StatusNoContent)
	return
}

func (s *Server) registerRoutes(h *Handler) {
	s.router.HandleFunc("GET /records", h.listHandler)
	s.router.HandleFunc("POST /records", h.createHandler)
	s.router.HandleFunc("DELETE /records/{id}", h.deleteHandler)
	s.router.HandleFunc("GET /records/{id}", h.retrieveHandler)
	s.router.HandleFunc("PUT /records/{id}", h.updateHandler)
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

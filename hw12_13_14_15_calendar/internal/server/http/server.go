package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SevaStopAll/golang-hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/SevaStopAll/golang-hw/hw12_13_14_15_16_calendar/internal/server/http/api"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	app    *app.App
}

func NewServer(app *app.App, host string, port int) *Server {
	s := &Server{app: app}
	router := s.setupRouter()

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupRouter() *mux.Router {
	router := mux.NewRouter()

	// API routes
	apiServer := api.NewServer(s.app)
	apiRouter := router.PathPrefix("/api").Subrouter()
	api.HandlerFromMux(apiServer, apiRouter)

	// Health check endpoint
	router.HandleFunc("/health", s.healthCheckHandler).Methods("GET")

	// OpenAPI specification
	router.HandleFunc("/openapi.yaml", s.openAPIHandler).Methods("GET")

	// Apply middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	return router
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) openAPIHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./api/openapi.yaml")
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запросов
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, duration)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

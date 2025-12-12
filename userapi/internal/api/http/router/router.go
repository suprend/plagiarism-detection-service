package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"userapi/internal/api/http/handler"
	"userapi/internal/application/usecase"
)

type Router struct {
	submitHandler    *handler.SubmitHandler
	reportsHandler   *handler.ReportsHandler
	wordcloudHandler *handler.WordcloudHandler
}

func NewRouter(submitUC *usecase.SubmitUseCase, reportsUC *usecase.ReportsUseCase, wcUC *usecase.WordcloudUseCase) *Router {
	return &Router{
		submitHandler:    handler.NewSubmitHandler(submitUC),
		reportsHandler:   handler.NewReportsHandler(reportsUC),
		wordcloudHandler: handler.NewWordcloudHandler(wcUC),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/works/", r.handleWorks)
	mux.HandleFunc("/wordcloud", r.wordcloudHandler.Handle)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/openapi.yaml", serveOpenAPI)
	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) handleWorks(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch {
	case strings.HasSuffix(path, "/submit"):
		r.submitHandler.Handle(w, req)
	case strings.HasSuffix(path, "/reports"):
		r.reportsHandler.Handle(w, req)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "not_found",
			"message": "unknown path",
		})
	}
}

func serveOpenAPI(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, req, "openapi.yaml")
}

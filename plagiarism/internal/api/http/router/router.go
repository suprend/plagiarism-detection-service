package router

import (
	"net/http"

	"plagiarism/internal/api/http/handler"
	"plagiarism/internal/application/usecase"
)

type Router struct {
	checkHandler   *handler.CheckHandler
	reportsHandler *handler.ReportsHandler
}

func NewRouter(checkUseCase usecase.CheckUseCase) *Router {
	return &Router{
		checkHandler:   handler.NewCheckHandler(checkUseCase),
		reportsHandler: handler.NewReportsHandler(checkUseCase),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/checks", r.checkHandler.Handle)
	mux.HandleFunc("/works/", r.reportsHandler.Handle)

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

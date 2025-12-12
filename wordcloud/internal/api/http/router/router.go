package router

import (
	"net/http"

	"wordcloud/internal/api/http/handler"
	"wordcloud/internal/application/usecase"
)

type Router struct {
	wordcloudHandler *handler.WordcloudHandler
}

func NewRouter(wc *usecase.WordcloudService) *Router {
	return &Router{
		wordcloudHandler: handler.NewWordcloudHandler(wc),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/wordcloud", r.wordcloudHandler.Handle)
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

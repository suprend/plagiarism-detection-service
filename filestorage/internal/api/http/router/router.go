package router

import (
	"net/http"

	"filestorage/internal/api/http/handler"
	"filestorage/internal/application/usecase"
)

type Router struct {
	submitHandler      *handler.SubmitHandler
	submissionsHandler *handler.SubmissionsHandler
	downloadHandler    *handler.DownloadHandler
}

func NewRouter(
	submitUseCase *usecase.SubmitUseCase,
	getSubmissionsUseCase *usecase.GetSubmissionsUseCase,
	downloadSubmissionUseCase *usecase.DownloadSubmissionUseCase,
) *Router {
	return &Router{
		submitHandler:      handler.NewSubmitHandler(submitUseCase),
		submissionsHandler: handler.NewSubmissionsHandler(getSubmissionsUseCase),
		downloadHandler:    handler.NewDownloadHandler(downloadSubmissionUseCase),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/submit", r.submitHandler.Handle)
	mux.HandleFunc("/submissions", r.submissionsHandler.Handle)
	mux.HandleFunc("/submissions/download", r.downloadHandler.Handle)

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

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

	return mux
}

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

	return mux
}

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
	return mux
}

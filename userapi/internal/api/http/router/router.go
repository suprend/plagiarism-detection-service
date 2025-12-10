package router

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"

	"userapi/internal/api/http/handler"
	"userapi/internal/application/usecase"
)

var (
	//go:embed swaggerui/*.html
	swaggerUI embed.FS

	swaggerHandler http.Handler
)

func init() {
	root, err := fs.Sub(swaggerUI, "swaggerui")
	if err != nil {
		panic("failed to load embedded swagger UI: " + err.Error())
	}
	swaggerHandler = http.FileServer(http.FS(root))
}

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
	mux.HandleFunc("/swagger", redirectSwaggerRoot)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerHandler))
	return mux
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

func redirectSwaggerRoot(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	target := "/swagger/"
	if req.URL.RawQuery != "" {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

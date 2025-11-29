package handler

import (
	"net/http"

	"userapi/internal/application/usecase"
)

type WordcloudHandler struct {
	useCase *usecase.WordcloudUseCase
}

func NewWordcloudHandler(uc *usecase.WordcloudUseCase) *WordcloudHandler {
	return &WordcloudHandler{useCase: uc}
}

func (h *WordcloudHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is allowed")
		return
	}

	submissionID := r.URL.Query().Get("submission_id")
	if submissionID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "submission_id is required")
		return
	}

	img, err := h.useCase.Generate(r.Context(), submissionID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "wordcloud_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(img)
}

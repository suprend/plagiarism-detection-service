package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"filestorage/internal/application/usecase"
)

type DownloadHandler struct {
	useCase *usecase.DownloadSubmissionUseCase
}

func NewDownloadHandler(useCase *usecase.DownloadSubmissionUseCase) *DownloadHandler {
	return &DownloadHandler{useCase: useCase}
}

func (h *DownloadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	submissionID := r.URL.Query().Get("submission_id")
	if submissionID == "" {
		respondValidationError(w, "submission_id query parameter is required")
		return
	}

	resp, err := h.useCase.Download(r.Context(), submissionID)
	if err != nil {
		log.Printf("download: submission_id=%s failed: %v", submissionID, err)
		respondError(w, err)
		return
	}
	defer resp.File.Close()

	w.Header().Set("Content-Type", resp.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, resp.Filename))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, resp.File); err != nil {
		log.Printf("download: submission_id=%s stream failed: %v", submissionID, err)
	}
}

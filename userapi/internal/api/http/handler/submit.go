package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"userapi/internal/application/dto"
	"userapi/internal/application/usecase"
	"userapi/internal/infrastructure/config"
)

var errFileTooLarge = errors.New("file too large")

type SubmitHandler struct {
	useCase *usecase.SubmitUseCase
}

func NewSubmitHandler(uc *usecase.SubmitUseCase) *SubmitHandler {
	return &SubmitHandler{useCase: uc}
}

func (h *SubmitHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is allowed")
		return
	}

	workID, ok := extractWorkID(r.URL.Path, "/submit")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_path", "expected /works/{work_id}/submit")
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_form", "expected multipart form data")
		return
	}

	maxUploadSize := config.MaxUploadSize()

	var (
		login       string
		fileData    []byte
		filename    string
		contentType string
	)

	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			writeError(w, http.StatusBadRequest, "invalid_form", "failed to read multipart body")
			return
		}

		switch part.FormName() {
		case "login":
			body, readErr := io.ReadAll(part)
			_ = part.Close()
			if readErr != nil {
				writeError(w, http.StatusBadRequest, "invalid_form", "failed to read login")
				return
			}
			login = string(body)
		case "file":
			data, ct, readErr := readFilePart(part, maxUploadSize)
			if readErr != nil {
				if errors.Is(readErr, errFileTooLarge) {
					writeError(w, http.StatusBadRequest, "validation_error", fmt.Sprintf("file exceeds max size %d bytes", maxUploadSize))
					return
				}
				writeError(w, http.StatusBadRequest, "invalid_form", "failed to read file")
				return
			}
			fileData = data
			filename = part.FileName()
			contentType = ct
		default:
			_ = part.Close()
		}
	}

	if login == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "login is required")
		return
	}

	if len(fileData) == 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req := dto.SubmitWorkRequest{
		WorkID:      workID,
		Login:       login,
		Data:        fileData,
		Filename:    filename,
		ContentType: contentType,
	}

	resp, err := h.useCase.Submit(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadGateway, "submit_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}

func readFilePart(part *multipart.Part, maxUploadSize int64) ([]byte, string, error) {
	defer part.Close()

	limited := io.LimitReader(part, maxUploadSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", err
	}

	if int64(len(data)) > maxUploadSize {
		return nil, "", errFileTooLarge
	}

	return data, part.Header.Get("Content-Type"), nil
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"filestorage/internal/application/dto"
	"filestorage/internal/application/usecase"
)

const maxUploadSize = 1 * 1024 * 1024 // 1MB

var errFileTooLarge = errors.New("file too large")

type SubmitHandler struct {
	submitUseCase *usecase.SubmitUseCase
}

func NewSubmitHandler(submitUseCase *usecase.SubmitUseCase) *SubmitHandler {
	return &SubmitHandler{
		submitUseCase: submitUseCase,
	}
}

func (h *SubmitHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondMethodNotAllowed(w, "only POST method is allowed")
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		respondValidationError(w, "invalid multipart form")
		return
	}

	var (
		assignmentID string
		login        string
		fileData     []byte
		filename     string
		contentType  string
	)

	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("submit: failed to read multipart part: %v", err)
			respondError(w, err)
			return
		}

		switch part.FormName() {
		case "assignment_id":
			body, readErr := io.ReadAll(part)
			_ = part.Close()
			if readErr != nil {
				respondValidationError(w, "failed to read assignment_id")
				return
			}
			assignmentID = string(body)
		case "login":
			body, readErr := io.ReadAll(part)
			_ = part.Close()
			if readErr != nil {
				respondValidationError(w, "failed to read login")
				return
			}
			login = string(body)
		case "file":
			data, readErr := readFilePart(part)
			if readErr != nil {
				if errors.Is(readErr, errFileTooLarge) {
					respondValidationError(w, fmt.Sprintf("file exceeds max size %d bytes", maxUploadSize))
					return
				}
				log.Printf("submit: failed to read file part: %v", readErr)
				respondValidationError(w, "failed to read file")
				return
			}
			fileData = data
			filename = part.FileName()
			contentType = part.Header.Get("Content-Type")
		default:
			_ = part.Close()
		}
	}

	if assignmentID == "" {
		respondValidationError(w, "assignment_id is required")
		return
	}

	if login == "" {
		respondValidationError(w, "login is required")
		return
	}

	if fileData == nil {
		respondValidationError(w, "file is required")
		return
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req := dto.SubmitRequest{
		AssignmentID: assignmentID,
		Login:        login,
		Data:         fileData,
		Filename:     filename,
		ContentType:  contentType,
	}

	resp, err := h.submitUseCase.Submit(r.Context(), req)
	if err != nil {
		log.Printf("submit: assignment_id=%s login=%s failed: %v", assignmentID, login, err)
		respondError(w, err)
		return
	}

	log.Printf("submit: assignment_id=%s login=%s submission_id=%s uploaded", assignmentID, login, resp.SubmissionID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"submission_id": resp.SubmissionID,
	})
}

func readFilePart(part io.ReadCloser) ([]byte, error) {
	defer part.Close()

	limited := io.LimitReader(part, maxUploadSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}

	if len(data) > maxUploadSize {
		return nil, errFileTooLarge
	}

	return data, nil
}

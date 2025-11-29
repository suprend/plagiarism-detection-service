package dto

import "plagiarism/internal/domain"

type StartCheckResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
}

type CheckStatusResponse struct {
	CheckReport domain.CheckReport `json:"report"`
}

type WorkReportsResponse struct {
	WorkID  string               `json:"work_id"`
	Reports []domain.CheckReport `json:"reports"`
}

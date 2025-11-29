package usecase

import (
	"context"

	"plagiarism/internal/application/dto"
	"plagiarism/internal/domain"
)

type CheckUseCase interface {
	StartCheck(ctx context.Context, submissionID, workID string) (*dto.StartCheckResponse, error)
	GetCheck(ctx context.Context, workID, submissionID string) (*dto.CheckStatusResponse, error)
	GetReportsByWork(ctx context.Context, workID string) (*dto.WorkReportsResponse, error)
}

var (
	_ domain.CheckReport
	_ dto.StartCheckResponse
)

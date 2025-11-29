package usecase

import (
	"context"
	"errors"

	"userapi/internal/application/dto"
)

var ErrReportNotFound = errors.New("report not found")

type ReportsProvider interface {
	GetReports(ctx context.Context, workID string) (*dto.WorkReportsResponse, error)
}

type ReportsUseCase struct {
	provider ReportsProvider
}

func NewReportsUseCase(provider ReportsProvider) *ReportsUseCase {
	return &ReportsUseCase{provider: provider}
}

func (uc *ReportsUseCase) GetByWork(ctx context.Context, workID string) (*dto.WorkReportsResponse, error) {
	return uc.provider.GetReports(ctx, workID)
}

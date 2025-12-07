package usecase

import (
	"context"
	"errors"

	"userapi/internal/application/dto"
	apperr "userapi/internal/common/errors"
	plagclient "userapi/internal/infrastructure/plagiarism"
)

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
	resp, err := uc.provider.GetReports(ctx, workID)
	if err != nil {
		if errors.Is(err, plagclient.ErrNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "report not found")
		}
		return nil, apperr.Wrap(err, apperr.CodeDownstream, "get reports failed")
	}
	return resp, nil
}

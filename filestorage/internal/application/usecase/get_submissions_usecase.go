package usecase

import (
	"context"

	"filestorage/internal/domain/entity"
	"filestorage/internal/domain/repository"
)

type GetSubmissionsUseCase struct {
	submissionRepo repository.SubmissionRepository
}

func NewGetSubmissionsUseCase(submissionRepo repository.SubmissionRepository) *GetSubmissionsUseCase {
	return &GetSubmissionsUseCase{
		submissionRepo: submissionRepo,
	}
}

func (uc *GetSubmissionsUseCase) GetByAssignmentID(ctx context.Context, assignmentID string) ([]*entity.Submission, error) {
	submissions, err := uc.submissionRepo.GetByAssignmentID(ctx, assignmentID)
	if err != nil {
		return nil, wrapDatabaseError(err, "failed to fetch submissions")
	}

	return submissions, nil
}

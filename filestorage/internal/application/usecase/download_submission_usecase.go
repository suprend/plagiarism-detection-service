package usecase

import (
	"context"
	"errors"
	"io"

	apperr "filestorage/internal/common/errors"
	"filestorage/internal/domain/repository"

	"github.com/aws/smithy-go"
	"github.com/google/uuid"
)

type DownloadSubmissionUseCase struct {
	submissionRepo repository.SubmissionRepository
	s3Repo         repository.S3Repository
}

type DownloadSubmissionResponse struct {
	File        io.ReadCloser
	Filename    string
	ContentType string
}

func NewDownloadSubmissionUseCase(
	submissionRepo repository.SubmissionRepository,
	s3Repo repository.S3Repository,
) *DownloadSubmissionUseCase {
	return &DownloadSubmissionUseCase{
		submissionRepo: submissionRepo,
		s3Repo:         s3Repo,
	}
}

func (uc *DownloadSubmissionUseCase) Download(ctx context.Context, submissionID string) (*DownloadSubmissionResponse, error) {
	if submissionID == "" {
		return nil, newValidationError("submission_id is required")
	}

	id, err := uuid.Parse(submissionID)
	if err != nil {
		return nil, newValidationError("invalid submission_id")
	}

	submission, err := uc.submissionRepo.GetByID(ctx, id)
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			return nil, err
		}
		return nil, wrapDatabaseError(err, "failed to get submission")
	}

	file, err := uc.s3Repo.GetFile(ctx, submission.SubmissionID.String())
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchKey" {
			return nil, wrapNotFoundError(err, "submission file not found")
		}
		return nil, wrapStorageError(err, "failed to get submission file")
	}

	return &DownloadSubmissionResponse{
		File:        file,
		Filename:    submission.SubmissionID.String(),
		ContentType: "application/octet-stream",
	}, nil
}

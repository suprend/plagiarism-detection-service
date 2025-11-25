package usecase

import (
	"context"
	"fmt"
	"log"

	"filestorage/internal/application/dto"
	"filestorage/internal/domain/repository"
)

type SubmitUseCase struct {
	submissionRepo repository.SubmissionRepository
	s3Repo         repository.S3Repository
}

func NewSubmitUseCase(
	submissionRepo repository.SubmissionRepository,
	s3Repo repository.S3Repository,
) *SubmitUseCase {
	return &SubmitUseCase{
		submissionRepo: submissionRepo,
		s3Repo:         s3Repo,
	}
}

func (uc *SubmitUseCase) Submit(ctx context.Context, req dto.SubmitRequest) (*dto.SubmitResponse, error) {
	// Создаем submission в БД в рамках транзакции
	submission, tx, err := uc.submissionRepo.CreateWithTx(ctx, req.AssignmentID, req.Login)
	if err != nil {
		return nil, wrapDatabaseError(err, "failed to create submission")
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Формируем ключ для S3: используем submission_id
	s3Key := submission.SubmissionID.String()

	// Загружаем файл в S3
	if err := uc.s3Repo.UploadFile(ctx, s3Key, req.Data, req.ContentType); err != nil {
		log.Printf("submit: submission_id=%s failed to upload to s3 key=%s: %v", submission.SubmissionID.String(), s3Key, err)
		return nil, wrapStorageError(err, "failed to upload file to storage")
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("submit: submission_id=%s commit failed, deleting s3 key=%s: %v", submission.SubmissionID.String(), s3Key, err)
		if delErr := uc.s3Repo.DeleteFile(ctx, s3Key); delErr != nil {
			cleanupErr := wrapStorageError(delErr, "failed to delete uploaded file after commit failure")
			return nil, wrapDatabaseError(fmt.Errorf("%v; cleanup: %v", err, cleanupErr), "failed to commit submission tx")
		}
		return nil, wrapDatabaseError(err, "failed to commit submission tx")
	}
	tx = nil

	return &dto.SubmitResponse{
		SubmissionID: submission.SubmissionID.String(),
	}, nil
}

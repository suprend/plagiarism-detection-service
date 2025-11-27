package repository

import (
	"context"

	"filestorage/internal/domain/entity"

	"github.com/google/uuid"
)

type SubmissionRepository interface {
	Create(ctx context.Context, assignmentID, authorID string) (*entity.Submission, error)

	CreateWithTx(ctx context.Context, assignmentID, authorID string) (*entity.Submission, Transaction, error)

	GetByID(ctx context.Context, submissionID uuid.UUID) (*entity.Submission, error)

	GetByAssignmentID(ctx context.Context, assignmentID string) ([]*entity.Submission, error)

	GetByAuthorID(ctx context.Context, authorID string) ([]*entity.Submission, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

package repository

import (
	"context"

	"filestorage/internal/domain/entity"

	"github.com/google/uuid"
)

// SubmissionRepository определяет интерфейс для работы с submissions в БД
type SubmissionRepository interface {
	// Create создает новую запись submission
	Create(ctx context.Context, assignmentID, authorID string) (*entity.Submission, error)

	// CreateWithTx создает submission в рамках транзакции, оставляя ее открытой для дальнейших действий
	CreateWithTx(ctx context.Context, assignmentID, authorID string) (*entity.Submission, Transaction, error)

	// GetByID получает submission по ID
	GetByID(ctx context.Context, submissionID uuid.UUID) (*entity.Submission, error)

	// GetByAssignmentID получает все submissions по assignment_id
	GetByAssignmentID(ctx context.Context, assignmentID string) ([]*entity.Submission, error)

	// GetByAuthorID получает все submissions по author_id
	GetByAuthorID(ctx context.Context, authorID string) ([]*entity.Submission, error)
}

// Transaction описывает контракт для управления транзакциями репозитория
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

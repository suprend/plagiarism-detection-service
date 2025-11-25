package postgres

import (
	"context"
	stdErrors "errors"

	apperr "filestorage/internal/common/errors"
	"filestorage/internal/domain/entity"
	"filestorage/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool    *pgxpool.Pool
	queries *Queries
}

// NewPostgresRepository создает новый экземпляр PostgreSQL репозитория
func NewPostgresRepository(pool *pgxpool.Pool) repository.SubmissionRepository {
	return &postgresRepository{
		pool:    pool,
		queries: New(pool),
	}
}

// toEntity конвертирует postgres.Submission в entity.Submission
func toEntity(pgSub Submission) *entity.Submission {
	return &entity.Submission{
		SubmissionID: pgSub.SubmissionID,
		AssignmentID: pgSub.AssignmentID,
		AuthorID:     pgSub.AuthorID,
		CreatedAt:    pgSub.CreatedAt,
	}
}

// toEntitySlice конвертирует слайс postgres.Submission в слайс entity.Submission
func toEntitySlice(pgSubs []Submission) []*entity.Submission {
	result := make([]*entity.Submission, 0, len(pgSubs))
	for _, pgSub := range pgSubs {
		result = append(result, toEntity(pgSub))
	}
	return result
}

func (r *postgresRepository) Create(ctx context.Context, assignmentID, authorID string) (*entity.Submission, error) {
	submission, tx, err := r.CreateWithTx(ctx, assignmentID, authorID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to commit submission tx")
	}

	return submission, nil
}

func (r *postgresRepository) CreateWithTx(ctx context.Context, assignmentID, authorID string) (*entity.Submission, repository.Transaction, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to begin submission tx")
	}

	queries := r.queries.WithTx(tx)
	pgSub, err := queries.CreateSubmission(ctx, CreateSubmissionParams{
		AssignmentID: assignmentID,
		AuthorID:     authorID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to create submission")
	}

	return toEntity(pgSub), &pgxTxWrapper{tx: tx}, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, submissionID uuid.UUID) (*entity.Submission, error) {
	pgSub, err := r.queries.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.Wrap(err, apperr.CodeNotFound, "submission not found")
		}
		return nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to get submission by id")
	}

	return toEntity(pgSub), nil
}

func (r *postgresRepository) GetByAssignmentID(ctx context.Context, assignmentID string) ([]*entity.Submission, error) {
	pgSubs, err := r.queries.GetSubmissionsByAssignmentID(ctx, assignmentID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to get submissions by assignment_id")
	}

	return toEntitySlice(pgSubs), nil
}

func (r *postgresRepository) GetByAuthorID(ctx context.Context, authorID string) ([]*entity.Submission, error) {
	pgSubs, err := r.queries.GetSubmissionsByAuthorID(ctx, authorID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDatabase, "failed to get submissions by author_id")
	}

	return toEntitySlice(pgSubs), nil
}

type pgxTxWrapper struct {
	tx pgx.Tx
}

func (t *pgxTxWrapper) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTxWrapper) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

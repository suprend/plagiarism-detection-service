package usecase

import (
	"context"
	"errors"
	"time"

	"plagiarism/internal/application/dto"
	apperr "plagiarism/internal/common/errors"
	"plagiarism/internal/domain"
	"plagiarism/internal/infrastructure/report"
)

type reportStore interface {
	Save(domain.CheckReport) error
	LoadBySubmissionID(workID, submissionID string) (domain.CheckReport, error)
	GetOverallByWork(workID string) ([]domain.CheckReport, error)
}

type worker interface {
	Enqueue(ctx context.Context, report domain.CheckReport) error
}

type CheckService struct {
	store  reportStore
	worker worker
}

var ErrCheckNotFound = apperr.New(apperr.CodeNotFound, "report not found")
var ErrWorkerUnavailable = apperr.New(apperr.CodeInternal, "worker not configured")

func NewCheckService(store reportStore, worker worker) *CheckService {
	return &CheckService{store: store, worker: worker}
}

func (s *CheckService) StartCheck(ctx context.Context, submissionID, workID string) (*dto.StartCheckResponse, error) {
	report := domain.CheckReport{
		WorkID:       workID,
		SubmissionID: submissionID,
		Status:       domain.CheckStatusPending,
		CreatedAt:    time.Now().UTC(),
	}

	if s.worker == nil {
		report.Status = domain.CheckStatusFailed
		report.Error = ErrWorkerUnavailable.Error()
		_ = s.store.Save(report)
		return nil, ErrWorkerUnavailable
	}

	if err := s.store.Save(report); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "save report failed")
	}

	if err := s.worker.Enqueue(ctx, report); err != nil {
		report.Status = domain.CheckStatusFailed
		report.Error = err.Error()
		_ = s.store.Save(report)
		return nil, apperr.Wrap(err, apperr.CodeInternal, "enqueue failed")
	}

	return &dto.StartCheckResponse{
		SubmissionID: submissionID,
		Status:       string(report.Status),
	}, nil
}

func (s *CheckService) GetCheck(ctx context.Context, workID, submissionID string) (*dto.CheckStatusResponse, error) {
	rep, err := s.store.LoadBySubmissionID(workID, submissionID)
	if err != nil {
		if errors.Is(err, report.ErrReportNotFound) {
			return nil, ErrCheckNotFound
		}
		return nil, apperr.Wrap(err, apperr.CodeInternal, "load report failed")
	}

	return &dto.CheckStatusResponse{CheckReport: rep}, nil
}

func (s *CheckService) GetReportsByWork(ctx context.Context, workID string) (*dto.WorkReportsResponse, error) {
	reports, err := s.store.GetOverallByWork(workID)
	if err != nil {
		if errors.Is(err, report.ErrReportNotFound) {
			return nil, ErrCheckNotFound
		}
		return nil, apperr.Wrap(err, apperr.CodeInternal, "get reports failed")
	}

	return &dto.WorkReportsResponse{
		WorkID:  workID,
		Reports: reports,
	}, nil
}

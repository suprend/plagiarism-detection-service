package usecase

import (
	"context"
	"fmt"

	"userapi/internal/application/dto"
)

type FilestorageUploader interface {
	UploadSubmission(ctx context.Context, assignmentID, login string, data []byte, filename, contentType string) (string, error)
}

type PlagiarismStarter interface {
	StartCheck(ctx context.Context, submissionID, workID string) (CheckStartResult, error)
}

type CheckStartResult struct {
	SubmissionID string
	Status       string
}

type SubmitUseCase struct {
	fs   FilestorageUploader
	plag PlagiarismStarter
}

func NewSubmitUseCase(fs FilestorageUploader, plag PlagiarismStarter) *SubmitUseCase {
	return &SubmitUseCase{
		fs:   fs,
		plag: plag,
	}
}

func (uc *SubmitUseCase) Submit(ctx context.Context, req dto.SubmitWorkRequest) (*dto.SubmitWorkResponse, error) {
	submissionID, err := uc.fs.UploadSubmission(ctx, req.WorkID, req.Login, req.Data, req.Filename, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("upload submission: %w", err)
	}

	check, err := uc.plag.StartCheck(ctx, submissionID, req.WorkID)
	if err != nil {
		return nil, fmt.Errorf("start plagiarism check: %w", err)
	}

	return &dto.SubmitWorkResponse{
		SubmissionID: check.SubmissionID,
		CheckStatus:  check.Status,
	}, nil
}

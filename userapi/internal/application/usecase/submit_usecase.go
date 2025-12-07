package usecase

import (
	"context"

	"userapi/internal/application/dto"
	apperr "userapi/internal/common/errors"
)

type FilestorageUploader interface {
	UploadSubmission(ctx context.Context, assignmentID, login string, data []byte, filename, contentType string) (string, error)
}

type PlagiarismStarter interface {
	StartCheck(ctx context.Context, submissionID, workID string) (dto.CheckStartResult, error)
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
		return nil, apperr.Wrap(err, apperr.CodeDownstream, "upload submission failed")
	}

	check, err := uc.plag.StartCheck(ctx, submissionID, req.WorkID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDownstream, "start plagiarism check failed")
	}

	return &dto.SubmitWorkResponse{
		SubmissionID: check.SubmissionID,
		CheckStatus:  check.Status,
	}, nil
}

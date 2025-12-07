package usecase

import (
	"context"

	apperr "userapi/internal/common/errors"
)

type wordcloudProvider interface {
	Get(ctx context.Context, submissionID string) ([]byte, error)
}

type WordcloudUseCase struct {
	provider wordcloudProvider
}

func NewWordcloudUseCase(provider wordcloudProvider) *WordcloudUseCase {
	return &WordcloudUseCase{provider: provider}
}

func (uc *WordcloudUseCase) Generate(ctx context.Context, submissionID string) ([]byte, error) {
	if submissionID == "" {
		return nil, apperr.New(apperr.CodeValidation, "submission_id is required")
	}
	img, err := uc.provider.Get(ctx, submissionID)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDownstream, "wordcloud request failed")
	}
	return img, nil
}

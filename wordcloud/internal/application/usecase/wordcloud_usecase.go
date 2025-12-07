package usecase

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"wordcloud/internal/common/errors"
)

type submissionDownloader interface {
	DownloadSubmission(ctx context.Context, submissionID string) ([]byte, error)
}

type wordcloudRenderer interface {
	Render(ctx context.Context, text string) ([]byte, error)
}

type WordcloudService struct {
	fs  submissionDownloader
	wc  wordcloudRenderer
	dir string
}

func NewWordcloudService(fs submissionDownloader, wc wordcloudRenderer, dir string) *WordcloudService {
	return &WordcloudService{fs: fs, wc: wc, dir: dir}
}

func (uc *WordcloudService) Generate(ctx context.Context, submissionID string) ([]byte, error) {
	if submissionID == "" {
		return nil, errors.New(errors.CodeValidation, "submission_id is required")
	}

	data, err := uc.fs.DownloadSubmission(ctx, submissionID)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeStorage, "download submission failed")
	}

	text := normalizeText(string(data))
	if text == "" {
		text = string(data)
	}
	if text == "" {
		return nil, errors.New(errors.CodeValidation, "submission is empty")
	}

	resp, err := uc.wc.Render(ctx, text)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeStorage, "render wordcloud failed")
	}

	if uc.dir != "" {
		if err := os.MkdirAll(uc.dir, 0o755); err == nil {
			path := filepath.Join(uc.dir, submissionID+".png")
			_ = os.WriteFile(path, resp, 0o644)
		}
	}
	return resp, nil
}

func normalizeText(raw string) string {
	if raw == "" {
		return ""
	}
	builder := strings.Builder{}
	builder.Grow(len(raw))
	for _, r := range raw {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '\n' {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	words := strings.Fields(strings.ToLower(builder.String()))
	filtered := make([]string, 0, len(words))
	for _, w := range words {
		filtered = append(filtered, w)
		if len(filtered) >= 5000 {
			break
		}
	}
	return strings.Join(filtered, " ")
}

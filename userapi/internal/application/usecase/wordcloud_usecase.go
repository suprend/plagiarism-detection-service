package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type submissionDownloader interface {
	DownloadSubmission(ctx context.Context, submissionID string) ([]byte, error)
}

type wordcloudRenderer interface {
	Render(ctx context.Context, text string) ([]byte, error)
}

type WordcloudUseCase struct {
	fs  submissionDownloader
	wc  wordcloudRenderer
	dir string
}

func NewWordcloudUseCase(fs submissionDownloader, wc wordcloudRenderer, dir string) *WordcloudUseCase {
	return &WordcloudUseCase{fs: fs, wc: wc, dir: dir}
}

func (uc *WordcloudUseCase) Generate(ctx context.Context, submissionID string) ([]byte, error) {
	if submissionID == "" {
		return nil, fmt.Errorf("submission_id is required")
	}

	data, err := uc.fs.DownloadSubmission(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("download submission: %w", err)
	}

	text := normalizeText(string(data))
	if text == "" {
		text = string(data)
	}

	resp, err := uc.wc.Render(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("render wordcloud: %w", err)
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

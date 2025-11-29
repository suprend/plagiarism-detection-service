package worker

import (
	"context"
	"fmt"
	"sync"

	"plagiarism/internal/domain"
	"plagiarism/internal/infrastructure/filestorage"
)

type Reporter interface {
	Save(domain.CheckReport) error
}

type FilestorageClient interface {
	ListSubmissions(ctx context.Context, assignmentID string) ([]filestorage.SubmissionMeta, error)
	DownloadSubmission(ctx context.Context, submissionID string) ([]byte, error)
}

type Worker struct {
	reporter  Reporter
	fs        FilestorageClient
	threshold float64
	onError   func(domain.CheckReport, error)

	tasks chan domain.CheckReport
	wg    sync.WaitGroup
}

func NewWorker(reporter Reporter, fs FilestorageClient, threshold float64, workers int, onError func(domain.CheckReport, error)) *Worker {
	if workers < 1 {
		workers = 1
	}

	w := &Worker{
		reporter:  reporter,
		fs:        fs,
		threshold: threshold,
		onError:   onError,
		tasks:     make(chan domain.CheckReport, 32),
	}

	w.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go w.loop()
	}
	return w
}

func (w *Worker) Close() {
	close(w.tasks)
	w.wg.Wait()
}

func (w *Worker) Enqueue(_ context.Context, report domain.CheckReport) error {
	if w.fs == nil {
		return fmt.Errorf("filestorage client not configured")
	}

	select {
	case w.tasks <- report:
		return nil
	default:
		return fmt.Errorf("worker queue is full")
	}
}

func (w *Worker) loop() {
	defer w.wg.Done()
	for report := range w.tasks {
		matches, selfAuthor, err := w.compareWithWork(report)
		if err != nil {
			report.Status = domain.CheckStatusFailed
			report.Error = err.Error()
		} else {
			report.Status = domain.CheckStatusDone
			report.Matches = matches
			report.AuthorID = selfAuthor
		}
		if err := w.reporter.Save(report); err != nil && w.onError != nil {
			w.onError(report, err)
		}
	}
}

func (w *Worker) compareWithWork(report domain.CheckReport) ([]domain.MatchResult, string, error) {
	ctx := context.Background()

	submissions, err := w.fs.ListSubmissions(ctx, report.WorkID)
	if err != nil {
		return nil, "", err
	}

	authors := make(map[string]string, len(submissions))
	for _, s := range submissions {
		authors[s.SubmissionID] = s.AuthorID
	}

	selfData, err := w.fs.DownloadSubmission(ctx, report.SubmissionID)
	if err != nil {
		return nil, "", err
	}

	matches := make([]domain.MatchResult, 0, len(submissions))
	selfAuthor := authors[report.SubmissionID]
	for _, sub := range submissions {
		if sub.SubmissionID == report.SubmissionID {
			continue
		}

		otherData, err := w.fs.DownloadSubmission(ctx, sub.SubmissionID)
		if err != nil {
			return nil, "", err
		}

		match := compareBytes(selfData, otherData, sub.SubmissionID, w.threshold)
		match.OtherAuthorID = authors[sub.SubmissionID]
		if match.Equal {
			matches = append(matches, match)
		}
	}

	return matches, selfAuthor, nil
}

func compareBytes(self, other []byte, otherID string, threshold float64) domain.MatchResult {
	matched := int64(0)
	total := int64(len(self))
	if len(other) > len(self) {
		total = int64(len(other))
	}

	minLen := len(self)
	if len(other) < minLen {
		minLen = len(other)
	}

	for i := 0; i < minLen; i++ {
		if self[i] == other[i] {
			matched++
		}
	}

	var ratio float64
	if total == 0 {
		ratio = 1
	} else {
		ratio = float64(matched) / float64(total)
	}
	equal := ratio >= threshold

	return domain.MatchResult{
		OtherSubmissionID: otherID,
		Equal:             equal,
		MatchedBytes:      matched,
		TotalBytes:        total,
		Similarity:        ratio,
		SelfSize:          int64(len(self)),
		OtherSize:         int64(len(other)),
	}
}

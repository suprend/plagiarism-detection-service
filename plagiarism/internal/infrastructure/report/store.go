package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"plagiarism/internal/domain"
)

var ErrReportNotFound = errors.New("report not found")

type FileReportStore struct {
	root string
	mu   sync.Mutex
}

func NewFileReportStore(root string) *FileReportStore {
	return &FileReportStore{root: root}
}

func (s *FileReportStore) Save(report domain.CheckReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	workDir := filepath.Join(s.root, sanitize(report.WorkID))
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(workDir, fmt.Sprintf("%s.json", sanitize(report.SubmissionID)))
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	return s.writeOverallLocked(workDir)
}

func (s *FileReportStore) LoadBySubmissionID(workID, submissionID string) (domain.CheckReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	target := filepath.Join(s.root, sanitize(workID), fmt.Sprintf("%s.json", sanitize(submissionID)))

	data, err := os.ReadFile(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.CheckReport{}, ErrReportNotFound
		}
		return domain.CheckReport{}, err
	}

	var report domain.CheckReport
	if err := json.Unmarshal(data, &report); err != nil {
		return domain.CheckReport{}, err
	}
	return report, nil
}

func (s *FileReportStore) GetOverallByWork(workID string) ([]domain.CheckReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	workDir := filepath.Join(s.root, sanitize(workID))
	data, err := os.ReadFile(filepath.Join(workDir, "overall.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}
	var payload struct {
		Reports []domain.CheckReport `json:"reports"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if len(payload.Reports) == 0 {
		return nil, ErrReportNotFound
	}
	return payload.Reports, nil
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	return s
}

func (s *FileReportStore) writeOverallLocked(workDir string) error {
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return err
	}
	var reports []domain.CheckReport
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") || entry.Name() == "overall.json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(workDir, entry.Name()))
		if err != nil {
			return err
		}
		var rep domain.CheckReport
		if err := json.Unmarshal(data, &rep); err != nil {
			return err
		}
		reports = append(reports, rep)
	}
	data, err := json.MarshalIndent(map[string]any{
		"reports": reports,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(workDir, "overall.json"), data, 0o644)
}

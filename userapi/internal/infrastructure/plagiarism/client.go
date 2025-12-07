package plagiarism

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

const checksPath = "/checks"

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type StartCheckResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
}

type WorkReportsResponse struct {
	WorkID  string        `json:"work_id"`
	Reports []CheckReport `json:"reports"`
}

type CheckReport struct {
	WorkID       string        `json:"work_id"`
	SubmissionID string        `json:"submission_id"`
	AuthorID     string        `json:"author_id"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	Error        string        `json:"error,omitempty"`
	Matches      []MatchResult `json:"matches"`
}

type MatchResult struct {
	OtherSubmissionID string  `json:"other_submission_id"`
	OtherAuthorID     string  `json:"other_author_id"`
	Equal             bool    `json:"equal"`
	MatchedBytes      int64   `json:"matched_bytes"`
	TotalBytes        int64   `json:"total_bytes"`
	Similarity        float64 `json:"similarity"`
	SelfSize          int64   `json:"self_size"`
	OtherSize         int64   `json:"other_size"`
}

func (c *Client) StartCheck(ctx context.Context, submissionID, workID string) (*StartCheckResponse, error) {
	payload := map[string]string{
		"submission_id": submissionID,
		"work_id":       workID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid plagiarism url: %w", err)
	}
	u.Path = checksPath

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("start check failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	var parsed StartCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if parsed.SubmissionID == "" {
		return nil, fmt.Errorf("start check failed: empty submission_id")
	}
	return &parsed, nil
}

func (c *Client) GetReports(ctx context.Context, workID string) (*WorkReportsResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid plagiarism url: %w", err)
	}
	u.Path = fmt.Sprintf("/works/%s/reports", url.PathEscape(workID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("get reports failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	var parsed WorkReportsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

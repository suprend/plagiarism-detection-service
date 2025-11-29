package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SubmissionMeta struct {
	SubmissionID string    `json:"submission_id"`
	AssignmentID string    `json:"assignment_id"`
	AuthorID     string    `json:"author_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ListSubmissions(ctx context.Context, assignmentID string) ([]SubmissionMeta, error) {
	u, err := url.Parse(c.baseURL + "/submissions")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("assignment_id", assignmentID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list submissions: status %d", resp.StatusCode)
	}

	var payload struct {
		Submissions []SubmissionMeta `json:"submissions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Submissions, nil
}

func (c *Client) DownloadSubmission(ctx context.Context, submissionID string) ([]byte, error) {
	u := fmt.Sprintf("%s/submissions/download?submission_id=%s", c.baseURL, url.QueryEscape(submissionID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download submission: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

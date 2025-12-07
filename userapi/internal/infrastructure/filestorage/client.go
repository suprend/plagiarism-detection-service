package filestorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

const submitPath = "/submit"
const downloadPath = "/submissions/download"

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) UploadSubmission(ctx context.Context, assignmentID, login string, data []byte, filename, contentType string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("assignment_id", assignmentID); err != nil {
		return "", err
	}
	if err := writer.WriteField("login", login); err != nil {
		return "", err
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid filestorage url: %w", err)
	}
	u.Path = submitPath

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if contentType != "" {
		req.Header.Set("X-File-Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("upload to filestorage failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	var payload struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.SubmissionID == "" {
		return "", fmt.Errorf("upload to filestorage failed: empty submission_id")
	}
	return payload.SubmissionID, nil
}

func (c *Client) DownloadSubmission(ctx context.Context, submissionID string) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid filestorage url: %w", err)
	}
	u.Path = downloadPath
	q := u.Query()
	q.Set("submission_id", submissionID)
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
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("download submission: status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

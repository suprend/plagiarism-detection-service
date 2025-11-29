package wordcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type requestPayload struct {
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Format     string `json:"format,omitempty"`
	FontFamily string `json:"fontFamily,omitempty"`
	Text       string `json:"text"`
}

func (c *Client) Render(ctx context.Context, text string) ([]byte, error) {
	payload := requestPayload{
		Width:      800,
		Height:     600,
		Format:     "png",
		FontFamily: "sans-serif",
		Text:       text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("wordcloud: status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

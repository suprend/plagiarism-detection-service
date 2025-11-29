package dto

import "time"

type MatchResult struct {
	OtherSubmissionID string  `json:"other_submission_id"`
	OtherAuthorID     string  `json:"other_author_id,omitempty"`
	Equal             bool    `json:"equal"`
	MatchedBytes      int64   `json:"matched_bytes"`
	TotalBytes        int64   `json:"total_bytes"`
	Similarity        float64 `json:"similarity"`
	SelfSize          int64   `json:"self_size"`
	OtherSize         int64   `json:"other_size"`
}

type CheckReport struct {
	WorkID       string        `json:"work_id"`
	SubmissionID string        `json:"submission_id"`
	AuthorID     string        `json:"author_id,omitempty"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	Error        string        `json:"error,omitempty"`
	Matches      []MatchResult `json:"matches"`
}

type WorkReportsResponse struct {
	WorkID  string        `json:"work_id"`
	Reports []CheckReport `json:"reports"`
}

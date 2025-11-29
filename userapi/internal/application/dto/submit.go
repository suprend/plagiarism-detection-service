package dto

type SubmitWorkRequest struct {
	WorkID      string
	Login       string
	Data        []byte
	Filename    string
	ContentType string
}

type SubmitWorkResponse struct {
	SubmissionID string `json:"submission_id"`
	CheckStatus  string `json:"check_status"`
}

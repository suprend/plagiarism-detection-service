package dto

type SubmitRequest struct {
	AssignmentID string
	Login        string
	Data         []byte
	Filename     string
	ContentType  string
}

type SubmitResponse struct {
	SubmissionID string
}

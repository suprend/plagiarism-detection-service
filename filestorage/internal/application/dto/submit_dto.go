package dto

// SubmitRequest представляет данные запроса на загрузку файла
type SubmitRequest struct {
	AssignmentID string
	Login        string
	Data         []byte
	Filename     string
	ContentType  string
}

// SubmitResponse представляет ответ на запрос загрузки
type SubmitResponse struct {
	SubmissionID string
}

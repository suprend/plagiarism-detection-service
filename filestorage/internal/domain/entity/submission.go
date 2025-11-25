package entity

import (
	"time"

	"github.com/google/uuid"
)

// Submission представляет сущность сдачи работы
type Submission struct {
	SubmissionID uuid.UUID
	AssignmentID string
	AuthorID     string
	CreatedAt    time.Time
}

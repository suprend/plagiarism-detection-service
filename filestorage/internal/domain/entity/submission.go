package entity

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	SubmissionID uuid.UUID
	AssignmentID string
	AuthorID     string
	CreatedAt    time.Time
}

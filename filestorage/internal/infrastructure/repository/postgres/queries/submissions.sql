-- name: CreateSubmission :one
INSERT INTO submissions (assignment_id, author_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetSubmissionByID :one
SELECT * FROM submissions
WHERE submission_id = $1;

-- name: GetSubmissionsByAssignmentID :many
SELECT * FROM submissions
WHERE assignment_id = $1
ORDER BY created_at DESC;

-- name: GetSubmissionsByAuthorID :many
SELECT * FROM submissions
WHERE author_id = $1
ORDER BY created_at DESC;


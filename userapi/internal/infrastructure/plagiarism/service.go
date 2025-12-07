package plagiarism

import (
	"context"

	"userapi/internal/application/dto"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) StartCheck(ctx context.Context, submissionID, workID string) (dto.CheckStartResult, error) {
	resp, err := s.client.StartCheck(ctx, submissionID, workID)
	if err != nil {
		return dto.CheckStartResult{}, err
	}
	return dto.CheckStartResult{
		SubmissionID: resp.SubmissionID,
		Status:       resp.Status,
	}, nil
}

func (s *Service) GetReports(ctx context.Context, workID string) (*dto.WorkReportsResponse, error) {
	resp, err := s.client.GetReports(ctx, workID)
	if err != nil {
		return nil, err
	}

	reports := make([]dto.CheckReport, 0, len(resp.Reports))
	for _, rep := range resp.Reports {
		matches := make([]dto.MatchResult, 0, len(rep.Matches))
		for _, m := range rep.Matches {
			matches = append(matches, dto.MatchResult{
				OtherSubmissionID: m.OtherSubmissionID,
				OtherAuthorID:     m.OtherAuthorID,
				Equal:             m.Equal,
				MatchedBytes:      m.MatchedBytes,
				TotalBytes:        m.TotalBytes,
				Similarity:        m.Similarity,
				SelfSize:          m.SelfSize,
				OtherSize:         m.OtherSize,
			})
		}
		reports = append(reports, dto.CheckReport{
			WorkID:       rep.WorkID,
			SubmissionID: rep.SubmissionID,
			AuthorID:     rep.AuthorID,
			Status:       rep.Status,
			CreatedAt:    rep.CreatedAt,
			Error:        rep.Error,
			Matches:      matches,
		})
	}

	return &dto.WorkReportsResponse{
		WorkID:  resp.WorkID,
		Reports: reports,
	}, nil
}

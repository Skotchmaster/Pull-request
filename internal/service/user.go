package service

import (
	"context"
	"pull_request/internal/models"
)

type UserRepository interface {
	SetIsActive(ctx context.Context, user_id string, is_active bool) (*models.User, error)
	GetReview(ctx context.Context, user_id string) ([]models.PullRequest, error)
}

type UserService struct {
	Repo UserRepository
}

func (s *UserService) SetIsActive(ctx context.Context, user_id string, is_active bool) (*models.User, error) {
	user, err := s.Repo.SetIsActive(ctx, user_id, is_active)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*models.PullRequestResponse, error) {

	prs, err := s.Repo.GetReview(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &models.PullRequestResponse{
		UserID:       userID,
		PullRequests: make([]models.PullRequestShort, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, models.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return resp, nil
}

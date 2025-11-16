package service

import (
	"context"
	"errors"
	"math/rand"
	"pull_request/internal/apierr"
	"pull_request/internal/models"
)

type PullRequestRepository interface {
	GetByID(ctx context.Context, prID string) (*models.PullRequest, error)
	CreateWithReviewers(ctx context.Context, pr *models.PullRequest, reviewerIDs []string) error

	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName string) ([]models.User, error)

	GetReviewers(ctx context.Context, prID string) ([]string, error)
	MarkMerged(ctx context.Context, prID string) error
	ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error
}

type PullRequestService struct {
	Repo PullRequestRepository
}

func toPullRequestDTO(pr *models.PullRequest, reviewers []string) *models.PullRequestDTO {
	return &models.PullRequestDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, prID, prName, authorID string) (*models.PullRequestDTO, error) {
	if _, err := s.Repo.GetByID(ctx, prID); err == nil {
		return nil, apierr.ErrPRExists
	} else if !errors.Is(err, apierr.ErrPRNotFound) {
		return nil, err
	}

	author, err := s.Repo.GetUserByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	members, err := s.Repo.GetActiveTeamMembers(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	reviewerIDs := make([]string, 0, 2)
	for _, u := range members {
		if u.UserID == authorID {
			continue
		}
		reviewerIDs = append(reviewerIDs, u.UserID)
		if len(reviewerIDs) == 2 {
			break
		}
	}

	pr := &models.PullRequest{
		PullRequestID:   prID,
		PullRequestName: prName,
		AuthorID:        authorID,
		Status:          models.PRStatus("OPEN"),
	}

	if err := s.Repo.CreateWithReviewers(ctx, pr, reviewerIDs); err != nil {
		return nil, err
	}

	return toPullRequestDTO(pr, reviewerIDs), nil
}

func (s *PullRequestService) MergePullRequest(ctx context.Context, prID string) (*models.PullRequestDTO, error) {
	pr, err := s.Repo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status != models.PRStatus("MERGED") {
		if err := s.Repo.MarkMerged(ctx, prID); err != nil {
			return nil, err
		}
		pr.Status = models.PRStatus("MERGED")
	}

	reviewers, err := s.Repo.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	return toPullRequestDTO(pr, reviewers), nil
}

func (s *PullRequestService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequestDTO, string, error) {
	pr, err := s.Repo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == models.PRStatus("MERGED") {
		return nil, "", apierr.ErrPRAlreadyMerged
	}

	oldUser, err := s.Repo.GetUserByID(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}

	reviewers, err := s.Repo.GetReviewers(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	isAssigned := false
	for _, id := range reviewers {
		if id == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", apierr.ErrReviewerNotAssigned
	}

	members, err := s.Repo.GetActiveTeamMembers(ctx, oldUser.TeamName)
	if err != nil {
		return nil, "", err
	}

	exclude := make(map[string]struct{}, len(reviewers)+2)
	exclude[oldUserID] = struct{}{}
	exclude[pr.AuthorID] = struct{}{}
	for _, id := range reviewers {
		exclude[id] = struct{}{}
	}

	var candidates []string
	for _, u := range members {
		if _, skip := exclude[u.UserID]; skip {
			continue
		}
		candidates = append(candidates, u.UserID)
	}

	if len(candidates) == 0 {
		return nil, "", apierr.ErrNoCandidate
	}

	idx := rand.Intn(len(candidates))
	newReviewerID := candidates[idx]

	if err := s.Repo.ReplaceReviewer(ctx, prID, oldUserID, newReviewerID); err != nil {
		return nil, "", err
	}

	updatedReviewers, err := s.Repo.GetReviewers(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	return toPullRequestDTO(pr, updatedReviewers), newReviewerID, nil
}

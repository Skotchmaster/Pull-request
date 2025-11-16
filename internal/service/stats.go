package service

import (
	"context"

	"pull_request/internal/models"
)

type StatsRepository interface {
	GetReviewerStats(ctx context.Context) ([]models.ReviewerStat, error)
}

type StatsService struct {
	Repo StatsRepository
}

func NewStatsService(repo StatsRepository) *StatsService {
	return &StatsService{Repo: repo}
}

func (s *StatsService) GetReviewerStats(ctx context.Context) ([]models.ReviewerStat, error) {
	return s.Repo.GetReviewerStats(ctx)
}

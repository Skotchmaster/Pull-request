package repo

import (
	"context"

	"pull_request/internal/models"
)

func (r *PullRequestRepo) GetReviewerStats(ctx context.Context) ([]models.ReviewerStat, error) {
	var stats []models.ReviewerStat

	if err := r.DB.WithContext(ctx).
		Table("pr_reviewers").
		Select("user_id, COUNT(*) AS assignments").
		Group("user_id").
		Order("assignments DESC").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

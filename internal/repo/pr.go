package repo

import (
	"context"
	"errors"
	"pull_request/internal/apierr"
	"pull_request/internal/models"
	"time"

	"gorm.io/gorm"
)

type PullRequestRepo struct {
	DB *gorm.DB
}

func (r *PullRequestRepo) GetByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	if err := r.DB.WithContext(ctx).
		Where("pr_id = ?", prID).
		First(&pr).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrPRNotFound
		}
		return nil, err
	}
	return &pr, nil
}

func (r *PullRequestRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PullRequestRepo) GetActiveTeamMembers(ctx context.Context, teamName string) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).
		Where("team_name = ? AND is_active = TRUE", teamName).
		Order("user_id ASC").
		Find(&users).Error; err != nil {

		return nil, err
	}
	return users, nil
}

func (r *PullRequestRepo) CreateWithReviewers(ctx context.Context, pr *models.PullRequest, reviewerIDs []string) error {
    return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(pr).Error; err != nil {
            return err
        }

        if len(reviewerIDs) == 0 {
            return nil
        }

        reviewers := make([]models.PRReviewer, 0, len(reviewerIDs))
        for _, id := range reviewerIDs {
            reviewers = append(reviewers, models.PRReviewer{
                PullRequestID: pr.PullRequestID,
                UserID:        id,
            })
        }

        if err := tx.Create(&reviewers).Error; err != nil {
            return err
        }

        return nil
    })
}


func (r *PullRequestRepo) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	var reviewerIDs []string

	if err := r.DB.WithContext(ctx).
		Table("pr_reviewers").
		Where("pr_id = ?", prID).
		Order("user_id ASC").
		Pluck("user_id", &reviewerIDs).Error; err != nil {

		return nil, err
	}

	return reviewerIDs, nil
}

func (r *PullRequestRepo) MarkMerged(ctx context.Context, prID string) error {
	return r.DB.WithContext(ctx).
		Model(&models.PullRequest{}).
		Where("pr_id = ?", prID).
		Updates(map[string]any{
			"status":    "MERGED",
			"merged_at": time.Now(),
		}).Error
}

func (r *PullRequestRepo) ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error {
    return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.
            Where("pr_id = ? AND user_id = ?", prID, oldUserID).
            Delete(&models.PRReviewer{}).Error; err != nil {
            return err
        }

        newReviewer := models.PRReviewer{
            PullRequestID: prID,
            UserID:        newUserID,
        }

        if err := tx.Create(&newReviewer).Error; err != nil {
            return err
        }

        return nil
    })
}

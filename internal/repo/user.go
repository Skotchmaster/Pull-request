package repo

import (
	"context"
	"errors"
	"pull_request/internal/apierr"
	"pull_request/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func(r *UserRepo) SetIsActive(ctx context.Context, user_id string, is_active bool) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrUserNotFound
		}
		return nil,err
	}

	user.IsActive = is_active

	if err := r.DB.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetReview(ctx context.Context, userID string) ([]models.PullRequest, error) {
    var user models.User
    if err := r.DB.WithContext(ctx).
        Where("user_id = ?", userID).
        First(&user).Error; err != nil {

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apierr.ErrUserNotFound
        }
        return nil, err
    }

    var prs []models.PullRequest
    if err := r.DB.WithContext(ctx).
        Model(&models.PullRequest{}).
        Joins("JOIN pr_reviewers prr ON prr.pr_id = pull_requests.pr_id").
        Where("prr.user_id = ?", userID).
        Order("pull_requests.created_at DESC").
        Find(&prs).Error; err != nil {

        return nil, err
    }

    return prs, nil
}

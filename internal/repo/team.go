package repo

import (
	"context"

	"pull_request/internal/models"
	"pull_request/internal/apierr"

	"gorm.io/gorm"
)

type TeamRepo struct {
	DB *gorm.DB
}

func (r *TeamRepo) CreateTeam(ctx context.Context, t models.TeamResp) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Team{}).
			Where("team_name = ?", t.TeamName).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return apierr.ErrTeamExists
		}

		if err := tx.Create(&models.Team{
			TeamName: t.TeamName,
		}).Error; err != nil {
			return err
		}

		for _, m := range t.Members {
			u := &models.User{
				UserID:   m.UserID,
				Username: m.Username,
				TeamName: t.TeamName,
				IsActive: m.IsActive,
			}

			if err := tx.Save(u).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TeamRepo) GetTeam(ctx context.Context, name string) (models.TeamResp, error) {

	var t models.Team
	if err := r.DB.WithContext(ctx).Where("team_name = ?", name).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.TeamResp{}, apierr.ErrTeamNotFound
		}
		return models.TeamResp{}, err
	}

	var users []models.User
	if err := r.DB.WithContext(ctx).Where("team_name = ?", name).Find(&users).Error; err != nil {
		return models.TeamResp{}, err
	}

	members := make([]models.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, models.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	return models.TeamResp{
		TeamName: t.TeamName,
		Members:  members,
	}, nil
}

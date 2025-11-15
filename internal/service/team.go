package service

import (
	"context"
	"pull_request/internal/models"
)



type TeamRepository interface {
	CreateTeam(ctx context.Context, t models.TeamResp) error
	GetTeam(ctx context.Context, name string) (models.TeamResp, error)
}

type TeamService struct {
	Repo TeamRepository
}

func (s *TeamService) CreateTeam(ctx context.Context, t models.TeamResp) (models.TeamResp, error) {
	if err := s.Repo.CreateTeam(ctx, t); err != nil {
		return models.TeamResp{}, err
	}
	return s.Repo.GetTeam(ctx, t.TeamName)
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (models.TeamResp, error) {
	return s.Repo.GetTeam(ctx, name)
}

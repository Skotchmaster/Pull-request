package service

import (
	"context"
)

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, t Team) error
	GetTeam(ctx context.Context, name string) (Team, error)
}

type TeamService struct {
	Repo TeamRepository
}

func (s *TeamService) CreateTeam(ctx context.Context, t Team) (Team, error) {
	if err := s.Repo.CreateTeam(ctx, t); err != nil {
		return Team{}, err
	}
	return s.Repo.GetTeam(ctx, t.TeamName)
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (Team, error) {
	return s.Repo.GetTeam(ctx, name)
}

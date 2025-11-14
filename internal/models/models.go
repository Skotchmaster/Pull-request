package models

import "time"

type Team struct {
	TeamName string `gorm:"column:team_name;primaryKey" json:"team_name"`
}

type User struct {
	UserID   string `gorm:"column:user_id;primaryKey" json:"user_id"`
	Username string `gorm:"column:username"            json:"username"`
	TeamName string `gorm:"column:team_name"           json:"team_name"`
	IsActive bool   `gorm:"column:is_active"           json:"is_active"`
}

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string     `gorm:"column:pr_id;primaryKey"   json:"pull_request_id"`
	PullRequestName string     `gorm:"column:pr_name"            json:"pull_request_name"`
	AuthorID        string     `gorm:"column:author_id"          json:"author_id"`
	Status          PRStatus   `gorm:"column:status"             json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at"         json:"createdAt"`
	MergedAt        *time.Time `gorm:"column:merged_at"          json:"mergedAt,omitempty"`

	AssignedReviewers []string `gorm:"-" json:"assigned_reviewers"`
}

type PRReviewer struct {
	PullRequestID string `gorm:"column:pr_id;primaryKey"  json:"pull_request_id"`
	UserID        string `gorm:"column:user_id;primaryKey" json:"user_id"`
}

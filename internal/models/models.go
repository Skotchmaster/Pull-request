package models

import "time"

type Team struct {
	TeamName string `gorm:"primaryKey"`
}

type User struct {
	UserID   string `gorm:"primaryKey"`
	Username string
	TeamName string
	IsActive bool
}

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PRID     string     `gorm:"primaryKey"`
	PRName   string
	AuthorID string
	Status   PRStatus

	CreatedAt time.Time
	MergedAt  *time.Time
}

type PrReviewer struct {
	PRID   string `gorm:"primaryKey"`
	UserID string `gorm:"primaryKey"`
}

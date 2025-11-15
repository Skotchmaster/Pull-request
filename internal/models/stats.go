package models

type ReviewerStat struct {
	UserID      string `json:"user_id"`
	Assignments int64  `json:"assignments"`
}
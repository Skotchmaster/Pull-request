package models

type PullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            PRStatus `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

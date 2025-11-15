package models

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          PRStatus `json:"status"`
}

type PullRequestResponse struct {
	UserID       string              `json:"user_id"`
	PullRequests []PullRequestShort  `json:"pull_requests"`
}
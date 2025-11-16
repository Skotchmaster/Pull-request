package apierr

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")

	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrPRNotFound           = errors.New("pr not found")
	ErrPRAlreadyMerged      = errors.New("pr already merged")
	ErrNoReviewersAvailable = errors.New("no reviewers available")

	ErrPRExists            = errors.New("pr already exists")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
)

type APIError struct {
	Error APIErrorBody `json:"error"`
}

type APIErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code, message string) APIError {
	return APIError{
		Error: APIErrorBody{
			Code:    code,
			Message: message,
		},
	}
}

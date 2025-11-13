package apierr

type APIError struct {
	Error APIErrorBody `json:"error"`
}

type APIErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code, msg string) APIError {
	return APIError{
		Error: APIErrorBody{
			Code:    code,
			Message: msg,
		},
	}
}

const (
	ErrCodeTeamExists   = "TEAM_EXISTS"
	ErrCodePRExists     = "PR_EXISTS"
	ErrCodePRMerged     = "PR_MERGED"
	ErrCodeNotAssigned  = "NOT_ASSIGNED"
	ErrCodeNoCandidate  = "NO_CANDIDATE"
	ErrCodeNotFound     = "NOT_FOUND"

	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
)

const (
	ErrMsgUnauthorized       = "missing or invalid token"
	ErrMsgAdminRequired      = "admin role required"
	ErrMsgResourceNotFound   = "resource not found"
)

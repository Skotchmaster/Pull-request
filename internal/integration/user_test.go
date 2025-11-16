package integration

import (
	"encoding/json"
	"net/http"
	"pull_request/internal/models"
	"testing"
)

func (env *testEnv) getReview(userID string) models.PullRequestResponse {
	env.t.Helper()

	path := "/users/getReview?user_id=" + userID

	resp, body := env.doJSON(http.MethodGet, path, nil)

	if resp.StatusCode != http.StatusOK {
		env.t.Fatalf("GET /users/getReview status=%d, body=%s", resp.StatusCode, string(body))
	}

	var out models.PullRequestResponse
	if err := json.Unmarshal(body, &out); err != nil {
		env.t.Fatalf("unmarshal review resp: %v, body=%s", err, string(body))
	}

	return out
}

func (env *testEnv) setUserIsActive(userID string, isActive bool) models.User {
	env.t.Helper()

	reqBody := map[string]any{
		"user_id":   userID,
		"is_active": isActive,
	}

	resp, body := env.doJSON(http.MethodPost, "/users/setIsActive", reqBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		env.t.Fatalf("POST /users/setIsActive status=%d, body=%s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		User models.User `json:"user"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal setIsActive resp: %v, body=%s", err, string(body))
	}

	return wrapper.User
}

func testUserGetReview_Success(t *testing.T, env *testEnv) {
	t.Helper()

	env.createTeam("team-int-user", []models.TeamMember{
		{UserID: "user-review-author-1", Username: "Author", IsActive: true},
		{UserID: "user-review-reviewer-1", Username: "Reviewer 1", IsActive: true},
	})

	pr := env.createPR("pr-user-review-1", "Integration PR", "user-review-author-1")

	reviewerID := pr.AssignedReviewers[0]

	resp := env.getReview(reviewerID)

	if resp.UserID != reviewerID {
		t.Fatalf("expected user_id=%q, got %q", reviewerID, resp.UserID)
	}
	if len(resp.PullRequests) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(resp.PullRequests))
	}
	if resp.PullRequests[0].PullRequestID != pr.PullRequestID {
		t.Fatalf("expected PR %q, got %q",
			pr.PullRequestID, resp.PullRequests[0].PullRequestID)
	}
}

func TestIntegration_UserGetReview(t *testing.T) {
	env := newTestEnv(t)
	testUserGetReview_Success(t, env)
}

func testUserSetIsActive_Success(t *testing.T, env *testEnv) {
	t.Helper()

	team := env.createTeam("team-int-active", []models.TeamMember{
		{UserID: "user-1", Username: "User 1", IsActive: true},
	})

	user := env.setUserIsActive("user-1", false)

	if user.UserID != "user-1" {
		t.Fatalf("expected user_id=user-1, got %q", user.UserID)
	}
	if user.TeamName != team.TeamName {
		t.Fatalf("expected team_name=%q, got %q", team.TeamName, user.TeamName)
	}
	if user.IsActive {
		t.Fatalf("expected is_active=false, got true")
	}
}

func TestIntegration_UserSetIsActive(t *testing.T) {
	env := newTestEnv(t)
	testUserSetIsActive_Success(t, env)
}

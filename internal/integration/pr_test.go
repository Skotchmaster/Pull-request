package integration

import (
	"encoding/json"
	"net/http"
	"pull_request/internal/models"
	"testing"
)

func (env *testEnv) createPR(prID, name, authorID string) models.PullRequestDTO {
	env.t.Helper()

	reqBody := map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": name,
		"author_id":         authorID,
	}

	resp, body := env.doJSON(http.MethodPost, "/pullRequest/create", reqBody)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		env.t.Fatalf("POST /pullRequest/create status=%d, body=%s",
			resp.StatusCode, string(body))
	}

	var wrapper struct {
		PR models.PullRequestDTO `json:"pr"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal pr resp: %v, body=%s", err, string(body))
	}

	return wrapper.PR
}

func (env *testEnv) mergePR(prID string) models.PullRequestDTO {
	env.t.Helper()

	reqBody := map[string]string{
		"pull_request_id": prID,
	}

	resp, body := env.doJSON(http.MethodPost, "/pullRequest/merge", reqBody)

	if resp.StatusCode != http.StatusOK {
		env.t.Fatalf("POST /pullRequest/merge status=%d, body=%s",
			resp.StatusCode, string(body))
	}

	var wrapper struct {
		PR models.PullRequestDTO `json:"pr"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal merge resp: %v, body=%s", err, string(body))
	}

	return wrapper.PR
}

func (env *testEnv) reassignReviewer(prID, oldUserID string) (models.PullRequestDTO, string) {
	env.t.Helper()

	reqBody := map[string]string{
		"pull_request_id": prID,
		"old_user_id":     oldUserID,
	}

	resp, body := env.doJSON(http.MethodPost, "/pullRequest/reassign", reqBody)

	if resp.StatusCode != http.StatusOK {
		env.t.Fatalf("POST /pullRequest/reassign status=%d, body=%s",
			resp.StatusCode, string(body))
	}

	var wrapper struct {
		PR         models.PullRequestDTO `json:"pr"`
		ReplacedBy string                `json:"replaced_by"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal reassign resp: %v, body=%s", err, string(body))
	}

	return wrapper.PR, wrapper.ReplacedBy
}

func testPullRequestCreate_Success(t *testing.T, env *testEnv) {
	t.Helper()

	env.createTeam("team-int-pr", []models.TeamMember{
		{UserID: "user-1", Username: "Author", IsActive: true},
		{UserID: "user-2", Username: "Reviewer 1", IsActive: true},
		{UserID: "user-3", Username: "Reviewer 2", IsActive: true},
	})

	pr := env.createPR("pr-1", "Integration PR", "user-1")

	if pr.PullRequestID != "pr-1" {
		t.Fatalf("expected pr_id=pr-1, got %q", pr.PullRequestID)
	}
	if pr.AuthorID != "user-1" {
		t.Fatalf("expected author_id=user-1, got %q", pr.AuthorID)
	}
	if len(pr.AssignedReviewers) == 0 {
		t.Fatalf("expected reviewers, got 0")
	}
}

func TestIntegration_PullRequestCreate(t *testing.T) {
	env := newTestEnv(t)
	testPullRequestCreate_Success(t, env)
}

func testPullRequestMerge_SuccessAndIdempotent(t *testing.T, env *testEnv) {
	t.Helper()

	env.createTeam("team-int-merge", []models.TeamMember{
		{UserID: "user-1", Username: "Author", IsActive: true},
		{UserID: "user-2", Username: "Reviewer 1", IsActive: true},
	})

	pr := env.createPR("pr-merge-1", "PR to merge", "user-1")

	if pr.Status != models.PROpen {
		t.Fatalf("expected new PR status OPEN, got %q", pr.Status)
	}

	merged := env.mergePR(pr.PullRequestID)
	if merged.Status != models.PRMerged {
		t.Fatalf("expected status MERGED after first merge, got %q", merged.Status)
	}

	mergedAgain := env.mergePR(pr.PullRequestID)
	if mergedAgain.Status != models.PRMerged {
		t.Fatalf("expected status MERGED after second (idempotent) merge, got %q", mergedAgain.Status)
	}
}

func TestIntegration_PullRequestMerge(t *testing.T) {
	env := newTestEnv(t)
	testPullRequestMerge_SuccessAndIdempotent(t, env)
}

func testPullRequestReassign_Success(t *testing.T, env *testEnv) {
	t.Helper()

	env.createTeam("team-int-reassign", []models.TeamMember{
		{UserID: "user-1", Username: "Author", IsActive: true},
		{UserID: "user-2", Username: "Reviewer 1", IsActive: true},
		{UserID: "user-3", Username: "Reviewer 2", IsActive: true},
		{UserID: "user-4", Username: "Reviewer 3", IsActive: true},
	})

	pr := env.createPR("pr-reassign-1", "PR to reassign", "user-1")

	if len(pr.AssignedReviewers) != 2 {
		t.Fatalf("expected 2 assigned reviewers, got %d", len(pr.AssignedReviewers))
	}

	oldReviewer := pr.AssignedReviewers[0]

	updated, replacedBy := env.reassignReviewer(pr.PullRequestID, oldReviewer)

	if replacedBy == "" {
		t.Fatalf("expected non-empty replacedBy")
	}
	if replacedBy == oldReviewer {
		t.Fatalf("expected replacedBy to be different from old reviewer (%s)", oldReviewer)
	}

	if len(updated.AssignedReviewers) != 2 {
		t.Fatalf("expected 2 reviewers after reassign, got %d", len(updated.AssignedReviewers))
	}

	foundOld := false
	foundNew := false
	for _, id := range updated.AssignedReviewers {
		if id == oldReviewer {
			foundOld = true
		}
		if id == replacedBy {
			foundNew = true
		}
	}

	if foundOld {
		t.Fatalf("expected old reviewer %s to be removed from PR reviewers", oldReviewer)
	}
	if !foundNew {
		t.Fatalf("expected new reviewer %s to be present in PR reviewers", replacedBy)
	}
}

func TestIntegration_PullRequestReassign(t *testing.T) {
	env := newTestEnv(t)
	testPullRequestReassign_Success(t, env)
}

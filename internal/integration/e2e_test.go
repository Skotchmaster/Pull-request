package integration

import (
	"testing"

	"pull_request/internal/models"
)

func TestE2E_PRLifecycleAndStats(t *testing.T) {
	env := newTestEnv(t)

	env.createTeam("team-e2e-lifecycle", []models.TeamMember{
		{UserID: "e2e-author-1", Username: "Author", IsActive: true},
		{UserID: "e2e-reviewer-1", Username: "Reviewer 1", IsActive: true},
		{UserID: "e2e-reviewer-2", Username: "Reviewer 2", IsActive: true},
		{UserID: "e2e-reviewer-3", Username: "Reviewer 3", IsActive: true},
	})

	pr := env.createPR("e2e-pr-1", "E2E PR lifecycle", "e2e-author-1")

	if pr.PullRequestID != "e2e-pr-1" {
		t.Fatalf("expected pull_request_id=e2e-pr-1, got %q", pr.PullRequestID)
	}
	if pr.AuthorID != "e2e-author-1" {
		t.Fatalf("expected author_id=e2e-author-1, got %q", pr.AuthorID)
	}
	if pr.Status != models.PROpen {
		t.Fatalf("expected status=OPEN, got %q", pr.Status)
	}
	if len(pr.AssignedReviewers) == 0 || len(pr.AssignedReviewers) > 2 {
		t.Fatalf("expected 1..2 assigned reviewers, got %d", len(pr.AssignedReviewers))
	}

	reviewerID := pr.AssignedReviewers[0]

	reviewResp := env.getReview(reviewerID)

	if reviewResp.UserID != reviewerID {
		t.Fatalf("expected review for user %q, got %q", reviewerID, reviewResp.UserID)
	}
	if len(reviewResp.PullRequests) == 0 {
		t.Fatalf("expected at least 1 PR in review list, got 0")
	}

	found := false
	for _, short := range reviewResp.PullRequests {
		if short.PullRequestID == pr.PullRequestID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected PR %q in review list, but it is missing", pr.PullRequestID)
	}

	oldReviewer := pr.AssignedReviewers[0]

	updatedPR, replacedBy := env.reassignReviewer(pr.PullRequestID, oldReviewer)

	if replacedBy == "" {
		t.Fatalf("expected non-empty replaced_by after reassign")
	}
	if replacedBy == oldReviewer {
		t.Fatalf("expected replaced_by != old reviewer (%s)", oldReviewer)
	}
	if len(updatedPR.AssignedReviewers) != 2 {
		t.Fatalf("expected 2 reviewers after reassign, got %d", len(updatedPR.AssignedReviewers))
	}

	foundOld := false
	foundNew := false
	for _, id := range updatedPR.AssignedReviewers {
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

	merged := env.mergePR(pr.PullRequestID)
	if merged.Status != models.PRMerged {
		t.Fatalf("expected status=MERGED after merge, got %q", merged.Status)
	}

	stats := env.getReviewerStats()
	if len(stats) == 0 {
		t.Fatalf("expected non-empty reviewer stats")
	}
}

func TestE2E_InactiveUsersAreNotAssigned(t *testing.T) {
	env := newTestEnv(t)

	env.createTeam("team-e2e-active", []models.TeamMember{
		{UserID: "e2e-author-2", Username: "Author", IsActive: true},
		{UserID: "e2e-reviewer-a", Username: "Reviewer A", IsActive: true},
		{UserID: "e2e-reviewer-b", Username: "Reviewer B", IsActive: true},
	})

	user := env.setUserIsActive("e2e-reviewer-b", false)
	if user.UserID != "e2e-reviewer-b" || user.IsActive {
		t.Fatalf("expected user %q to be inactive after setIsActive, got is_active=%v",
			"e2e-reviewer-b", user.IsActive)
	}

	pr := env.createPR("e2e-pr-2", "E2E PR inactive users", "e2e-author-2")

	if pr.Status != models.PROpen {
		t.Fatalf("expected new PR status OPEN, got %q", pr.Status)
	}
	if len(pr.AssignedReviewers) == 0 {
		t.Fatalf("expected at least 1 reviewer, got 0")
	}

	for _, id := range pr.AssignedReviewers {
		if id == "e2e-reviewer-b" {
			t.Fatalf("expected inactive reviewer %q not to be assigned", id)
		}
	}
}

package integration

import (
	"encoding/json"
	"net/http"
	"pull_request/internal/models"
	"testing"
)

func (env *testEnv) getReviewerStats() []models.ReviewerStat {
	env.t.Helper()

	resp, body := env.doJSON(http.MethodGet, "/stats/reviewers", nil)

	if resp.StatusCode != http.StatusOK {
		env.t.Fatalf("GET /stats/reviewers status=%d, body=%s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		ReviewerStats []models.ReviewerStat `json:"reviewer_stats"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal stats resp: %v, body=%s", err, string(body))
	}

	return wrapper.ReviewerStats
}

func testStatsReviewers_Success(t *testing.T, env *testEnv) {
	t.Helper()

	env.createTeam("team-int-stats", []models.TeamMember{
		{UserID: "stats-user-1", Username: "Author", IsActive: true},
		{UserID: "stats-user-2", Username: "Reviewer", IsActive: true},
	})

	env.createPR("pr-stats-1", "Integration PR for stats", "stats-user-1")

	stats := env.getReviewerStats()
	if len(stats) == 0 {
		t.Fatalf("expected non-empty stats")
	}
}

func TestIntegration_StatsReviewers(t *testing.T) {
	env := newTestEnv(t)
	testStatsReviewers_Success(t, env)
}

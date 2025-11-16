package integration

import (
	"encoding/json"
	"net/http"
	"pull_request/internal/models"
	"testing"
)

func (env *testEnv) createTeam(teamName string, members []models.TeamMember) models.TeamResp {
	env.t.Helper()

	reqBody := models.TeamResp{
		TeamName: teamName,
		Members:  members,
	}

	resp, body := env.doJSON(http.MethodPost, "/team/add", reqBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		env.t.Fatalf("POST /team/add status=%d, body=%s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Team models.TeamResp `json:"team"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		env.t.Fatalf("unmarshal team resp: %v, body=%s", err, string(body))
	}

	return wrapper.Team
}

func (env *testEnv) getTeam(teamName string) models.TeamResp {
	env.t.Helper()

	path := "/team/get?team_name=" + teamName

	resp, body := env.doJSON(http.MethodGet, path, nil)

	if resp.StatusCode != http.StatusOK {
		env.t.Fatalf("GET /team/get status=%d, body=%s", resp.StatusCode, string(body))
	}

	var out models.TeamResp
	if err := json.Unmarshal(body, &out); err != nil {
		env.t.Fatalf("unmarshal team get resp: %v, body=%s", err, string(body))
	}

	return out
}

func testTeamAdd_Success(t *testing.T, env *testEnv) {
	t.Helper()

	resp := env.createTeam("team-int-1", []models.TeamMember{
		{UserID: "user-1", Username: "Author", IsActive: true},
		{UserID: "user-2", Username: "Reviewer", IsActive: true},
	})

	if resp.TeamName != "team-int-1" {
		t.Fatalf("expected team_name=team-int-1, got %q", resp.TeamName)
	}
	if len(resp.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(resp.Members))
	}
}

func TestIntegration_TeamAdd(t *testing.T) {
	env := newTestEnv(t)
	testTeamAdd_Success(t, env)
}

func testTeamGet_Success(t *testing.T, env *testEnv) {
	t.Helper()

	created := env.createTeam("team-int-get", []models.TeamMember{
		{UserID: "user-1", Username: "User 1", IsActive: true},
		{UserID: "user-2", Username: "User 2", IsActive: false},
	})

	got := env.getTeam(created.TeamName)

	if got.TeamName != created.TeamName {
		t.Fatalf("expected team_name=%q, got %q", created.TeamName, got.TeamName)
	}
	if len(got.Members) != len(created.Members) {
		t.Fatalf("expected %d members, got %d", len(created.Members), len(got.Members))
	}
}

func TestIntegration_TeamGet(t *testing.T) {
	env := newTestEnv(t)
	testTeamGet_Success(t, env)
}

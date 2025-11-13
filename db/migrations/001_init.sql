CREATE TABLE teams (team_name TEXT PRIMARY KEY);
CREATE TABLE users (
  user_id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE RESTRICT,
  is_active BOOLEAN NOT NULL
);
CREATE TABLE pull_requests (
  pr_id TEXT PRIMARY KEY,
  pr_name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  status TEXT NOT NULL CHECK (status IN ('OPEN','MERGED')),
  created_at TIMESTAMPTZ DEFAULT now(),
  merged_at TIMESTAMPTZ
);
CREATE TABLE pr_reviewers (
  pr_id   TEXT NOT NULL REFERENCES pull_requests(pr_id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  PRIMARY KEY (pr_id, user_id)
);
CREATE INDEX idx_users_team ON users(team_name);
CREATE INDEX idx_pr_status  ON pull_requests(status);
CREATE INDEX idx_rev_user   ON pr_reviewers(user_id);

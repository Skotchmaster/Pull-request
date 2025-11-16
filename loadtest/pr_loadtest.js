import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 20,
  duration: '30s',
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';

const TEAM_COUNT = 20;
const USERS_PER_TEAM = 10;

const TEAMS = Array.from({ length: TEAM_COUNT }, (_, teamIndex) => {
  const teamNumber = teamIndex + 1;
  const teamName = `loadtest-team-${teamNumber}`;

  const members = Array.from({ length: USERS_PER_TEAM }, (_, userIndex) => {
    const globalUserIndex = teamIndex * USERS_PER_TEAM + userIndex + 1;
    return {
      user_id: `u${globalUserIndex}`,
      username: `User ${globalUserIndex}`,
      is_active: true,
    };
  });

  return {
    team_name: teamName,
    members,
  };
});

export function setup() {
  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  for (const team of TEAMS) {
    const payload = JSON.stringify(team);
    const res = http.post(`${BASE_URL}/team/add`, payload, params);

    check(res, {
      [`team ${team.team_name} created or already exists`]: (r) =>
        (r.status >= 200 && r.status < 300) || r.status === 409,
    });
  }
}

export default function () {
  const jsonHeaders = { headers: { 'Content-Type': 'application/json' } };

  const team = TEAMS[Math.floor(Math.random() * TEAMS.length)];
  const author =
    team.members[Math.floor(Math.random() * team.members.length)];

  const prId = `pr-${team.team_name}-${__VU}-${__ITER}`;

  const createRes = http.post(
    `${BASE_URL}/pullRequest/create`,
    JSON.stringify({
      pull_request_id: prId,
      pull_request_name: `Load test PR ${prId}`,
      author_id: author.user_id,
    }),
    jsonHeaders,
  );

  check(createRes, {
    'create PR: 2xx или 409 (PR_EXISTS)': (r) =>
      (r.status >= 200 && r.status < 300) || r.status === 409,
  });

  const reviewRes = http.get(
    `${BASE_URL}/users/getReview?user_id=${author.user_id}`,
  );

  check(reviewRes, {
    'getReview 200': (r) => r.status === 200,
  });

  const statsRes = http.get(`${BASE_URL}/stats/reviewers`);

  check(statsRes, {
    'stats 200': (r) => r.status === 200,
  });

  sleep(0.1);
}

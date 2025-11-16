package handlers

type Handlers struct {
	PRHandler    PullRequestHandler
	TeamHandler  TeamHandler
	UserHandler  UserHandler
	StatsHandler StatsHandler
}

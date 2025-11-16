package http

import (
	"pull_request/internal/handlers"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	Handlers handlers.Handlers
}

func Register(e *echo.Echo, d *Deps) {
	e.GET("/health/live", func(c echo.Context) error { return c.NoContent(200) })

	e.POST("/team/add", d.Handlers.TeamHandler.CreateTeam)
	e.GET("/team/get", d.Handlers.TeamHandler.GetTeam)

	e.GET("/users/getReview", d.Handlers.UserHandler.GetUserReviews)
	e.POST("/users/setIsActive", d.Handlers.UserHandler.SetUserIsActive)

	e.POST("/pullRequest/create", d.Handlers.PRHandler.CreatePullRequest)
	e.POST("/pullRequest/merge", d.Handlers.PRHandler.MergePullRequest)
	e.POST("/pullRequest/reassign", d.Handlers.PRHandler.ReassignReviewer)

	e.GET("/stats/reviewers", d.Handlers.StatsHandler.GetReviewerStats)
}

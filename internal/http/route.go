package http

import (
	"pull_request/internal/auth"
	"pull_request/internal/handlers"
	"pull_request/internal/logging"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	JWTSecret []byte
	Middleware auth.Middleware
	Handlers handlers.Handlers
}

func Register(e *echo.Echo, d *Deps, logger logging.Logger) {
	e.GET("/health/live", func(c echo.Context) error { return c.NoContent(200) })

	e.POST("/team/add", d.Handlers.TeamHandler.CreateTeam)

	auth := e.Group("", d.Middleware.RequireLogin(logger))
	auth.GET("/team/get", d.Handlers.TeamHandler.GetTeam)
	auth.GET("/users/getReview", d.Handlers.UserHandler.GetUserReviews)

	admin := e.Group("", d.Middleware.RequireAdmin(logger))
	admin.POST("/users/setIsActive", d.Handlers.UserHandler.SetUserIsActive)
	admin.POST("/pullRequest/create", d.Handlers.PRHandler.CreatePullRequest)
	admin.POST("/pullRequest/merge", d.Handlers.PRHandler.MergePullRequest)
	admin.POST("/pullRequest/reassign", d.Handlers.PRHandler.ReassignReviewer)
}
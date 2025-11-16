package handlers

import (
	"net/http"
	"pull_request/internal/apierr"
	"pull_request/internal/logging"
	"pull_request/internal/service"

	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	Service *service.StatsService
	Logger  logging.Logger
}

func (h *StatsHandler) GetReviewerStats(c echo.Context) error {
	ctx := c.Request().Context()

	stats, err := h.Service.GetReviewerStats(ctx)
	if err != nil {
		h.Logger.Errorf("failed to get reviewer stats: %v", err)
		return c.JSON(http.StatusInternalServerError,
			apierr.NewAPIError("INTERNAL", "internal error"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviewer_stats": stats,
	})
}

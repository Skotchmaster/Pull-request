package handlers

import (
	"errors"
	"net/http"

	"pull_request/internal/apierr"
	"pull_request/internal/logging"
	"pull_request/internal/service"

	"github.com/labstack/echo/v4"
)

type TeamHandler struct {
	Service *service.TeamService
	Logger  logging.Logger
}

func (h *TeamHandler) CreateTeam(c echo.Context) error {
	var req service.Team

	if err := c.Bind(&req); err != nil {
		h.Logger.Warnf("failed to bind /team/add request: %v", err)
		return c.JSON(http.StatusBadRequest, apierr.NewAPIError("BAD_REQUEST", "invalid request body"))
	}

	if req.TeamName == "" {
		h.Logger.Warnf("team_name is required")
		return c.JSON(http.StatusBadRequest, apierr.NewAPIError("BAD_REQUEST", "team_name is required"))
	}

	team, err := h.Service.CreateTeam(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apierr.ErrTeamExists):
			h.Logger.Warnf("team_name already exists")
			return c.JSON(http.StatusBadRequest, apierr.NewAPIError("TEAM_EXISTS", "team_name already exists"))
		default:
			h.Logger.Errorf("failed to create team: %v", err)
			return c.JSON(http.StatusInternalServerError, apierr.NewAPIError("INTERNAL", "internal error"))
		}
	}

	return c.JSON(http.StatusCreated, map[string]any{"team": team})
}

func (h *TeamHandler) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		h.Logger.Warnf("team_name is required")
		return c.JSON(http.StatusBadRequest, apierr.NewAPIError("BAD_REQUEST", "team_name is required"))
	}

	team, err := h.Service.GetTeam(c.Request().Context(), teamName)
	if err != nil {
		switch {
		case errors.Is(err, apierr.ErrTeamNotFound):
			h.Logger.Warnf("team not found")
			return c.JSON(http.StatusNotFound, apierr.NewAPIError("NOT_FOUND", "team not found"))
		default:
			h.Logger.Errorf("failed to get team %s: %v", teamName, err)
			return c.JSON(http.StatusInternalServerError, apierr.NewAPIError("INTERNAL", "internal error"))
		}
	}

	return c.JSON(http.StatusOK, team)
}

package handlers

import (
	"errors"
	"net/http"
	"pull_request/internal/apierr"
	"pull_request/internal/logging"
	"pull_request/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *service.UserService
	logger logging.Logger
}

func (h *UserHandler) SetUserIsActive(c echo.Context) error {
    ctx := c.Request().Context()

	type setUserIsActiveRequest struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}


    var req setUserIsActiveRequest
    if err := c.Bind(&req); err != nil {
        h.logger.Warnf("failed to bind request: %v", err)
        return c.JSON(http.StatusBadRequest,
            apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
    }

    if req.UserID == "" {
        h.logger.Warnf("user_id is required")
        return c.JSON(http.StatusBadRequest,
            apierr.NewAPIError("BAD_REQUEST", "user_id is required"))
    }

    user, err := h.service.SetIsActive(ctx, req.UserID, req.IsActive)
    if err != nil {
        if errors.Is(err, apierr.ErrUserNotFound){
            h.logger.Warnf("user not found: %s", req.UserID)
            return c.JSON(http.StatusNotFound, apierr.NewAPIError("NOT_FOUND", "user not found"))
		}
        h.logger.Errorf("internal server error: %v", err)
        return c.JSON(http.StatusInternalServerError, apierr.NewAPIError("INTERNAL", "internal error"))
    }

    return c.JSON(http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) GetUserReviews(c echo.Context) error {
    ctx := c.Request().Context()

    userID := c.QueryParam("user_id")
    if userID == "" {
        h.logger.Warnf("user_id is required")
        return c.JSON(http.StatusBadRequest,
            apierr.NewAPIError("BAD_REQUEST", "user_id is required"))
    }

    resp, err := h.service.GetReview(ctx, userID)
    if err != nil {
        if errors.Is(err, apierr.ErrUserNotFound) {
            h.logger.Warnf("user not found: %s", userID)
            return c.JSON(http.StatusNotFound, apierr.NewAPIError("NOT_FOUND", "user not found"))
        }
        h.logger.Errorf("internal server error: %v", err)
        return c.JSON(http.StatusInternalServerError, apierr.NewAPIError("INTERNAL", "internal error"))
    }

    return c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"errors"
	"net/http"
	"pull_request/internal/apierr"
	"pull_request/internal/logging"
	"pull_request/internal/service"

	"github.com/labstack/echo/v4"
)

type PullRequestHandler struct {
	service *service.PullRequestService
	logger  logging.Logger
}

func NewPullRequestHandler(s *service.PullRequestService, logger logging.Logger) *PullRequestHandler {
	return &PullRequestHandler{
		service: s,
		logger:  logger,
	}
}

type createPullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

func (h *PullRequestHandler) CreatePullRequest(c echo.Context) error {
	ctx := c.Request().Context()

	var req createPullRequestRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warnf("failed to bind PR create request: %v", err)
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		h.logger.Warnf("invalid PR create payload: %#v", req)
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
	}

	pr, err := h.service.CreatePullRequest(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, apierr.ErrUserNotFound):
			h.logger.Warnf("author or team not found %s: %v", req.PullRequestID, err)
			return c.JSON(http.StatusNotFound,
				apierr.NewAPIError("NOT_FOUND", "author or team not found"))
		case errors.Is(err, apierr.ErrPRExists):
			h.logger.Warnf("pull request already exists: %s", req.PullRequestID)
			return c.JSON(http.StatusConflict,
				apierr.NewAPIError("PR_EXISTS", "PR id already exists"))
		default:
			h.logger.Errorf("failed to create pull request: %v", err)
			return c.JSON(http.StatusInternalServerError,
				apierr.NewAPIError("INTERNAL", "internal error"))
		}
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"pr": pr,
	})
}

type mergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (h *PullRequestHandler) MergePullRequest(c echo.Context) error {
	ctx := c.Request().Context()

	var req mergePullRequestRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warnf("failed to bind PR merge request: %v", err)
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
	}

	if req.PullRequestID == "" {
		h.logger.Warnf("pull_request_id is required for merge")
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "pull_request_id is required"))
	}

	pr, err := h.service.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, apierr.ErrPRNotFound):
			h.logger.Warnf("pull request not found: %s", req.PullRequestID)
			return c.JSON(http.StatusNotFound,
				apierr.NewAPIError("NOT_FOUND", "pull request not found"))
		default:
			h.logger.Errorf("failed to merge pull request: %v", err)
			return c.JSON(http.StatusInternalServerError,
				apierr.NewAPIError("INTERNAL", "internal error"))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"pr": pr,
	})
}

type reassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (h *PullRequestHandler) ReassignReviewer(c echo.Context) error {
	ctx := c.Request().Context()

	var req reassignReviewerRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warnf("failed to bind PR reassign request: %v", err)
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		h.logger.Warnf("invalid PR reassign payload: %#v", req)
		return c.JSON(http.StatusBadRequest,
			apierr.NewAPIError("BAD_REQUEST", "invalid payload"))
	}

	pr, replacedBy, err := h.service.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, apierr.ErrPRNotFound),
			errors.Is(err, apierr.ErrUserNotFound):
			h.logger.Warnf("PR or user not found (pr_id=%s, user_id=%s): %v", req.PullRequestID, req.OldUserID, err)
			return c.JSON(http.StatusNotFound,
				apierr.NewAPIError("NOT_FOUND", "pr or user not found"))
		case errors.Is(err, apierr.ErrPRAlreadyMerged):
			h.logger.Warnf("attempt to reassign reviewer on merged PR: %s", req.PullRequestID)
			return c.JSON(http.StatusConflict,
				apierr.NewAPIError("PR_MERGED", "cannot reassign on merged PR"))
		case errors.Is(err, apierr.ErrReviewerNotAssigned):
			h.logger.Warnf("user is not assigned as reviewer (pr_id=%s, user_id=%s)", req.PullRequestID, req.OldUserID)
			return c.JSON(http.StatusConflict,
				apierr.NewAPIError("NOT_ASSIGNED", "reviewer is not assigned to this PR"))
		case errors.Is(err, apierr.ErrNoCandidate):
			h.logger.Warnf("no replacement candidate for reviewer (pr_id=%s, user_id=%s)", req.PullRequestID, req.OldUserID)
			return c.JSON(http.StatusConflict,
				apierr.NewAPIError("NO_CANDIDATE", "no active replacement candidate in team"))
		default:
			h.logger.Errorf("failed to reassign reviewer: %v", err)
			return c.JSON(http.StatusInternalServerError,
				apierr.NewAPIError("INTERNAL", "internal error"))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"pr":          pr,
		"replaced_by": replacedBy,
	})
}

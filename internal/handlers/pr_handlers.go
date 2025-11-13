package handlers

import (
	"pull_request/internal/logging"

	"gorm.io/gorm"
)

type PRHandler struct {
	logger logging.Logger
	db     *gorm.DB
}
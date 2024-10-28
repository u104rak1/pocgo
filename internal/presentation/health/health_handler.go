package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type HealthHandler struct {
	DB *bun.DB
}

func NewHealthHandler(db *bun.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

type HealthResponseBody struct {
	Status   string `json:"status" example:"healthy"`
	Database string `json:"database" example:"connected"`
}

func (h *HealthHandler) Run(ctx echo.Context) error {
	if err := h.DB.Ping(); err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, HealthResponseBody{
			Status:   "unhealthy",
			Database: "disconnected",
		})
	}

	return ctx.JSON(http.StatusOK, HealthResponseBody{
		Status:   "healthy",
		Database: "connected",
	})
}

package controller

import (
	"{{.ModulePath}}/internal/config"

	"github.com/linkeunid/ligo"
)

// HealthController handles health check requests.
type HealthController struct {
	cfg *config.Config
}

// NewHealthController creates a new health controller.
func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{cfg: cfg}
}

// Routes registers all routes for the health controller.
func (c *HealthController) Routes(r ligo.Router) {
	r.Handle("GET", "/health", c.Check)
}

// Check handles GET /health
func (c *HealthController) Check(ctx ligo.Context) error {
	return ctx.JSON(200, map[string]string{
		"status":  "ok",
		"version": c.cfg.Version,
	})
}

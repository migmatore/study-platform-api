package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	httpLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/migmatore/study-platform-api/pkg/logger"
)

type Deps struct {
	AuthUseCase AuthUseCase
}

type Handler struct {
	app    *fiber.App
	logger logger.Logger

	auth *AuthHandler
}

func New(deps Deps) *Handler {
	return &Handler{auth: NewUserHandler(deps.AuthUseCase)}
}

func (h *Handler) Init(ctx context.Context) *fiber.App {
	h.app = fiber.New()

	h.app.Use(cors.New())
	h.app.Use(httpLog.New())
	h.app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)

		return c.Next()
	})

	api := h.app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/signin", h.auth.Signin)
	auth.Post("/signup", h.auth.Signup)

	return h.app
}

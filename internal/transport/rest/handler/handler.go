package handler

import (
	"context"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	httpLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/pkg/jwt"
)

type Deps struct {
	AuthUseCase      AuthUseCase
	ClassroomUseCase ClassroomUseCase
	LessonUseCase    LessonUseCase
	StudentUseCase   StudentUseCase
}

type Handler struct {
	config *config.Config
	app    *fiber.App

	auth      *AuthHandler
	classroom *ClassroomHandler
	lesson    *LessonHandler
	student   *StudentHandler
}

func New(config *config.Config, deps Deps) *Handler {
	return &Handler{
		config:    config,
		auth:      NewUserHandler(deps.AuthUseCase),
		classroom: NewClassroomHandler(deps.ClassroomUseCase, deps.LessonUseCase),
		lesson:    NewLessonHundler(deps.LessonUseCase),
		student:   NewStudentsHandler(deps.StudentUseCase),
	}
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
	auth.Post("/refresh", h.auth.Refresh)

	v1.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(h.config.Server.JwtSecretKey)},
		ContextKey:   "jwt",
		ErrorHandler: jwt.JwtError,
	}))
	classrooms := v1.Group("/classrooms")
	classrooms.Get("/", h.classroom.All)
	classrooms.Post("/", h.classroom.Create)
	classrooms.Get("/:id/lessons", h.classroom.Lessons)
	classrooms.Get("/:id/lessons/current", h.classroom.CurrentLesson)
	classrooms.Post("/:id/lessons", h.classroom.CreateLesson)
	classrooms.Put("/:id/lessons", h.classroom.UpdateLesson)

	classrooms.Get("/:id/students", h.classroom.Students)

	lessons := v1.Group("/lessons")
	lessons.Get("/:id", h.lesson.ById)

	students := v1.Group("/students")
	students.Get("/", h.student.Students)

	return h.app
}

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
	UserUseCase      UserUseCase
	ClassroomUseCase ClassroomUseCase
	LessonUseCase    LessonUseCase
	StudentUseCase   StudentUseCase
	TeacherUseCase   TeacherUseCase
}

type Handler struct {
	config *config.Config
	app    *fiber.App

	auth      *AuthHandler
	user      *UserHandler
	classroom *ClassroomHandler
	lesson    *LessonHandler
	student   *StudentHandler
	teacher   *TeacherHandler
}

func New(config *config.Config, deps Deps) *Handler {
	return &Handler{
		config:    config,
		auth:      NewAuthHandler(deps.AuthUseCase),
		user:      NewUserHandler(deps.UserUseCase),
		classroom: NewClassroomHandler(deps.ClassroomUseCase, deps.LessonUseCase),
		lesson:    NewLessonHandler(deps.LessonUseCase),
		student:   NewStudentsHandler(deps.StudentUseCase),
		teacher:   NewTeacherHandler(deps.TeacherUseCase),
	}
}

func (h *Handler) Init(ctx context.Context) *fiber.App {
	h.app = fiber.New()

	h.app.Use(cors.New(cors.Config{
		AllowOrigins: "https://learnflow.ru",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	//h.app.Use(cors.New())
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
	users := v1.Group("/users")
	users.Get("/profile", h.user.Profile)
	users.Put("/profile", h.user.UpdateProfile)

	classrooms := v1.Group("/classrooms")
	classrooms.Get("/", h.classroom.All)
	classrooms.Post("/", h.classroom.Create)
	classrooms.Delete("/:id", h.classroom.Delete)
	classrooms.Get("/:id/lessons", h.classroom.Lessons)
	classrooms.Get("/:id/lessons/current", h.classroom.CurrentLesson)
	classrooms.Post("/:id/lessons", h.classroom.CreateLesson)
	classrooms.Put("/:id/lessons", h.classroom.UpdateLesson)

	classrooms.Get("/:id/students", h.classroom.Students)

	lessons := v1.Group("/lessons")
	lessons.Get("/:id", h.lesson.ById)
	lessons.Delete("/:id", h.lesson.Delete)

	students := v1.Group("/students")
	students.Get("/", h.student.Students)
	students.Post("/", h.student.Create)
	students.Delete("/:id", h.student.Delete)

	teachers := v1.Group("/teachers")
	teachers.Get("/", h.teacher.All)
	teachers.Post("/", h.teacher.Create)
	teachers.Delete("/:id", h.teacher.Delete)

	return h.app
}

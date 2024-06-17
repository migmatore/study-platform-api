package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	httpLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"log"
)

type AuthUseCase interface {
	Auth(ctx context.Context, req core.UserAuthRequest) (core.TokenMetadata, error)
	Refresh(ctx context.Context, req core.UserTokenRefreshRequest) (core.UserAuthResponse, error)
}

type ClassroomUseCase interface {
	Students(ctx context.Context, metadata core.TokenMetadata, classroomId int) ([]core.StudentResponse, error)
}

type HandlerDeps struct {
	AuthUseCase      AuthUseCase
	ClassroomUseCase ClassroomUseCase
}

type Handler struct {
	config *config.Config
	app    *fiber.App

	hub              *Hub
	authUseCase      AuthUseCase
	classroomUseCase ClassroomUseCase
}

func NewHandler(config *config.Config, deps HandlerDeps) *Handler {
	return &Handler{
		config:           config,
		hub:              NewHub(),
		authUseCase:      deps.AuthUseCase,
		classroomUseCase: deps.ClassroomUseCase,
	}
}

func (h *Handler) Init() *fiber.App {
	h.app = fiber.New()

	h.app.Use(cors.New(cors.Config{
		AllowOrigins: "https://learnflow.ru",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	h.app.Use(httpLog.New())

	go h.hub.Run()

	h.app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		// TODO REFACTOR THIS SHIT!!!!!!
		_, m, err := conn.ReadMessage()
		if err != nil {
			log.Println("error ", err)
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()

			return
		}

		req := struct {
			Type  MessageType `json:"type"`
			Token *string     `json:"token,omitempty"`
		}{}

		if err := json.Unmarshal(m, &req); err != nil {
			log.Println("error ", err)
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()

			return
		}

		if req.Type != AuthRequest {
			log.Println("error user must request auth")
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()

			return
		}

		if req.Token == nil || *req.Token == "" {
			log.Println("error user must provide token")
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()

			return
		}

		metadata, err := h.authUseCase.Auth(context.Background(), core.UserAuthRequest{Token: *req.Token})
		if err != nil {
			log.Println("error ", err)

			if errors.Is(err, apperrors.ExpiredToken) {
				conn.WriteJSON(map[string]interface{}{
					"type":       ErrorResp,
					"error_type": ExpiredTokenError,
				})
			}

			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()

			return
		}

		client := NewClient(
			ClientArgs{
				hub:      h.hub,
				conn:     conn,
				userId:   metadata.UserId,
				userRole: core.RoleType(metadata.Role),
			},
			ClientDeps{classroomUseCase: h.classroomUseCase},
		)

		h.hub.register <- client

		go client.writePump()

		client.readPump()
	}))

	return h.app
}

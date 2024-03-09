package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"os"
	"os/signal"
)

type Server struct {
	app    *fiber.App
	addr   string
	logger logger.Logger
}

func NewWebsocketServer(addr string, app *fiber.App, logger logger.Logger) *Server {
	return &Server{
		app:    app,
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start() {
	if err := s.app.Listen(s.addr); err != nil {
		s.logger.Fatalf("Websocket server is not running! Reason: %v", err)
	}
}

func (s *Server) StartWithGracefulShutdown() {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.app.Shutdown(); err != nil {
			s.logger.Errorf("Websocket server is not shutting down! Reason: %v", err)
		}

		s.logger.Info("Websocket server has successfully shut down!")

		close(idleConnsClosed)
	}()

	if err := s.app.Listen(s.addr); err != nil {
		s.logger.Fatalf("Websocket server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

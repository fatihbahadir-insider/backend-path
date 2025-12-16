package server

import (
	"backend-path/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	App *fiber.App
	DB *gorm.DB
}

func New(app *fiber.App, db *gorm.DB) *Server {
	return &Server{
		App: app,
		DB: db,
	}
}

func (s *Server) Start(port string) error {
	s.gracefulShutdown()

	utils.Logger.Info("Server is running on port " + port)
	return s.App.Listen(":" + port)
}

func (s *Server) gracefulShutdown() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		utils.Logger.Info("Shutting down server...")

		if err := s.App.Shutdown(); err != nil {
			utils.Logger.Error("Error shutting down server: " + err.Error())
		}

		s.cleanup()
	}()
}

func (s *Server) cleanup() {
	if s.DB != nil {
		if sqlDB, err := s.DB.DB(); err == nil {
			sqlDB.Close()
			utils.Logger.Info("Database connection closed")
		}
	}

	utils.Logger.Info("Server shutdown complete")
}
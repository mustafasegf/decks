package api

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mustafasegf/decks/internal/logger"
	"github.com/mustafasegf/decks/pkg/decks"
)

func (s *Server) SetupRouter() {
	s.Router.Use(logger.MiddleWare())
	s.Router.Use(recover.New())

	decksRepo := decks.NewRepo(s.DB)
	decksService := decks.NewService(decksRepo)
	decksHandler := decks.NewHandler(decksService)

	s.Router.Get("/decks/create", decksHandler.CreateNewDeck)
	s.Router.Get("/decks/open", decksHandler.OpenDeck)
	s.Router.Get("/decks/draw", decksHandler.DrawCard)
}

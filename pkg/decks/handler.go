package decks

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// @Summary Create a new deck
// @Description Endpoint for creating a brand new deck.
// @Tags /decks
// @Produce json
// @Param cards query string false "cards"
// @Param shuffled query string false "shuffled"
// @Success 200 {object} decks.CreateNewDeckResponse
// @Failure 400 {object} decks.ErrorResponse
// @Failure 500 {object} decks.ErrorResponse
// @Router /decks/create [get]
func (h *Handler) CreateNewDeck(c *fiber.Ctx) error {
	cards := c.Query("cards")
	shuffled := c.Query("shuffled") != ""
	res, status, err := h.Service.CreateNewDeck(cards, shuffled)

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(res)
}

// @Summary Open a new deck
// @Description Endpoint for returning whole deck.
// @Tags /decks
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} decks.OpenDeckResponse
// @Failure 400 {object} decks.ErrorResponse
// @Failure 500 {object} decks.ErrorResponse
// @Router /decks/open [get]
func (h *Handler) OpenDeck(c *fiber.Ctx) error {
	id := c.Query("id")
	res, status, err := h.Service.OpenDeck(id)

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(res)
}

// @Summary Draw a card from a deck
// @Description Endpoint for draw cards.
// @Tags /decks
// @Produce json
// @Param id query string true "id"
// @Param count query string true "count"
// @Success 200 {object} decks.DrawCardResponse
// @Failure 400 {object} decks.ErrorResponse
// @Failure 500 {object} decks.ErrorResponse
// @Router /decks/draw [get]
func (h *Handler) DrawCard(c *fiber.Ctx) error {
	countStr := c.Query("count", "0")
	id := c.Query("id")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "count must be a number",
		})
	}

	res, status, err := h.Service.DrawCard(id, count)

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(res)
}

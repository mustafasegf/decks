package decks

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mustafasegf/decks/pkg/models"
	"golang.org/x/exp/slices"
)

type Service interface {
	CreateNewDeck(cardsStr string, shuffled bool) (res CreateNewDeckResponse, status int, err error)
	OpenDeck(deckID string) (res OpenDeckResponse, status int, err error)
	DrawCard(deckID string, count int) (res DrawCardResponse, status int, err error)
}

type service struct {
	Repo Repo
}

func NewService(repo Repo) *service {
	svc := &service{
		Repo: repo,
	}

	return svc
}

func (s *service) CreateNewDeck(cardsStr string, shuffled bool) (res CreateNewDeckResponse, status int, err error) {
	status = http.StatusOK
	uuid := uuid.New()
	cards, err := s.generateDeck(cardsStr, shuffled)

	if err != nil {
		status = http.StatusBadRequest
		return
	}
	card := models.Deck{
		UUID:     uuid.String(),
		Shuffled: shuffled,
		Cards:    strings.Join(cards, ","),
	}

	err = s.Repo.CreateCard(card)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	res = CreateNewDeckResponse{
		DeckID:    uuid.String(),
		Shuffled:  shuffled,
		Remaining: len(cards),
	}

	return
}

func (s *service) OpenDeck(deckID string) (res OpenDeckResponse, status int, err error) {
	status = http.StatusOK
	deck, err := s.Repo.GetDeck(deckID)

	if err != nil {
		return
	}

	cards, err := s.convertToCards(deck.Cards)
	if err != nil {
		status = http.StatusBadRequest
		return
	}

	res = OpenDeckResponse{
		DeckID:    deck.UUID,
		Shuffled:  deck.Shuffled,
		Remaining: len(cards),
		Cards:     cards,
	}

	return
}

func (s *service) DrawCard(deckID string, count int) (res DrawCardResponse, status int, err error) {
	status = http.StatusOK
	if count <= 0 {
		status = fiber.StatusBadRequest
		err = fmt.Errorf("count must be greater than 0")
		return
	}

	deck, err := s.Repo.GetDeck(deckID)
	if err != nil {
		return
	}

	cards, cardsStr, err := s.removeCards(deck.Cards, count)
	if err != nil {
		status = fiber.StatusBadRequest
	}
	deck.Cards = cardsStr

	err = s.Repo.UpdateDeck(deck)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	res = DrawCardResponse{
		Cards: cards,
	}
	return
}

func (s *service) removeCards(cardsRaw string, count int) (cards []Card, cardsStr string, err error) {
	stringslice := strings.Split(cardsRaw, ",")
	if len(stringslice) < count {
		count = len(stringslice)
	}
	cardsStr = strings.Join(stringslice[count:], ",")

	cardsRaw = strings.Join(stringslice[:count], ",")
	cards, err = s.convertToCards(cardsRaw)

	return
}

func (s *service) convertToCards(cardsStr string) (cards []Card, err error) {
	suitMap := map[string]string{"S": "SPADES", "H": "HEARTS", "D": "DIAMONDS", "C": "CLUBS"}
	rankMap := map[string]string{"A": "ACE", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9", "10": "10", "J": "JACK", "Q": "QUEEN", "K": "KING"}

	cardsStrSlice := strings.Split(cardsStr, ",")
	cards = make([]Card, 0, len(cardsStrSlice))
	if len(cardsStr) == 0 {
		return
	}

	for _, cardStr := range cardsStrSlice {
		n := len(cardStr)
		if n < 2 {
			err = fmt.Errorf("card must be at least 2 characters long")
			return
		}
		val, ok := rankMap[cardStr[:n-1]]
		if !ok {
			err = fmt.Errorf("card %s is not valid rank", cardStr)
			return
		}
		suit, ok := suitMap[cardStr[n-1:]]
		if !ok {
			err = fmt.Errorf("card %s is not valid suit", cardStr)
			return
		}
		card := Card{
			Value: val,
			Suit:  suit,
			Code:  cardStr,
		}
		cards = append(cards, card)
	}
	return
}

func (s *service) generateDeck(cardsStr string, shuffled bool) (cards []string, err error) {
	cardsStrSlice := []string{"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS", "AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD", "AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC", "AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH"}

	cards = cardsStrSlice
	if cardsStr != "" {
		cards = strings.Split(cardsStr, ",")
		for _, card := range cards {
			card := strings.TrimSpace(card)
			if !slices.Contains(cardsStrSlice, card) {
				err = fmt.Errorf("card %s is not valid", card)
				return
			}
		}
	}

	if shuffled {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	return
}

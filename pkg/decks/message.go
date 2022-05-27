package decks

type CreateNewDeckResponse struct {
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type OpenDeckResponse struct {
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []Card `json:"cards"`
}

type DrawCardResponse struct {
	Cards []Card `json:"cards"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Expected something"`
}

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

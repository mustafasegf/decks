package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Deck struct {
	UUID      string    `json:"uuid"`
	ID        int       `json:"id"`
	Shuffled  bool      `json:"shuffled"`
	Cards     string    `json:"cards"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

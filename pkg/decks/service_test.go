package decks

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCardService(t *testing.T) {
	repo := NewRepo(db)
	svc := NewService(repo)

	t.Run("no parameters", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("", false)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusOK, expectedStr(http.StatusOK, status))
		assert.NotNil(t, res, expectedStr("not nil", res))

		assert.NotNil(t, res.DeckID, expectedStr("UUID", res.DeckID))
		assert.Equal(t, res.Remaining, 52, expectedStr("UUID", res.Remaining))
		assert.Equal(t, res.Shuffled, false, expectedStr(false, res.Shuffled))

		uuid, err := uuid.Parse(res.DeckID)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.NotNil(t, uuid, expectedStr("not nil", uuid))

	})

	t.Run("with shuffled", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("", true)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusOK, expectedStr(http.StatusOK, status))
		assert.NotNil(t, res, expectedStr("not nil", res))

		assert.NotNil(t, res.DeckID, expectedStr("UUID", res.DeckID))
		assert.Equal(t, res.Remaining, 52, expectedStr("UUID", res.Remaining))
		assert.Equal(t, res.Shuffled, true, expectedStr(true, res.Shuffled))

		uuid, err := uuid.Parse(res.DeckID)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.NotNil(t, uuid, expectedStr("not nil", uuid))
	})

	t.Run("with custom card", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("AS,2D,3C,KH", false)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusOK, expectedStr(http.StatusOK, status))
		assert.NotNil(t, res, expectedStr("not nil", res))

		assert.NotNil(t, res.DeckID, expectedStr("UUID", res.DeckID))
		assert.Equal(t, res.Remaining, 4, expectedStr("UUID", res.Remaining))
		assert.Equal(t, res.Shuffled, false, expectedStr(false, res.Shuffled))

		uuid, err := uuid.Parse(res.DeckID)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.NotNil(t, uuid, expectedStr("not nil", uuid))
	})

	t.Run("with custom card and shuffled", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("AS,2D,3C,KH", true)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusOK, expectedStr(http.StatusOK, status))
		assert.NotNil(t, res, expectedStr("not nil", res))

		assert.NotNil(t, res.DeckID, expectedStr("UUID", res.DeckID))
		assert.Equal(t, res.Remaining, 4, expectedStr("UUID", res.Remaining))
		assert.Equal(t, res.Shuffled, true, expectedStr(true, res.Shuffled))

		uuid, err := uuid.Parse(res.DeckID)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.NotNil(t, uuid, expectedStr("not nil", uuid))
	})

	t.Run("with invalid card", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("random", false)
		assert.Error(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusBadRequest, expectedStr(http.StatusBadRequest, status))
		assert.Empty(t, res, expectedStr("nil", res))
	})

	t.Run("with spaces in between comma", func(t *testing.T) {
		res, status, err := svc.CreateNewDeck("AS, 2D , 3C       ,  KH ", false)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, status, http.StatusOK, expectedStr(http.StatusOK, status))
		assert.NotNil(t, res, expectedStr("not nil", res))

		assert.NotNil(t, res.DeckID, expectedStr("UUID", res.DeckID))
		assert.Equal(t, res.Remaining, 4, expectedStr("UUID", res.Remaining))
		assert.Equal(t, res.Shuffled, false, expectedStr(false, res.Shuffled))

		uuid, err := uuid.Parse(res.DeckID)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.NotNil(t, uuid, expectedStr("not nil", uuid))
	})
}

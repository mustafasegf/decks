package decks

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mustafasegf/decks/pkg/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		zap.L().Fatal("Could not connect to docker", zap.Error(err))
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.2-alpine3.15",
		Env:        []string{"POSTGRES_PASSWORD=secret"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		zap.L().Fatal("Could not start resource", zap.Error(err))
	}

	port := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("host=localhost user=postgres password=secret dbname=postgres port=%s sslmode=disable", port)

	zap.L().Info("Connecting to database on: ", zap.String("dsn", dsn))

	resource.Expire(120)
	pool.MaxWait = 120 * time.Second

	if err := pool.Retry(func() error {
		var err error
		ctx := context.Background()
		db, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return db.Ping(ctx)

	}); err != nil {
		zap.L().Fatal("Could not connect to database", zap.Error(err))
	}

	dbURL := fmt.Sprintf("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", port)
	migr, err := migrate.New(
		"file://../../migrations",
		dbURL,
	)

	if err != nil {
		zap.L().Fatal("Could not create migration", zap.Error(err))
	}

	if err := migr.Up(); err != nil {
		zap.L().Fatal("Could not migrate", zap.Error(err))
	}
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		zap.L().Fatal("Could not purge resource", zap.Error(err))
	}

	os.Exit(code)
}

func TestCardRepo(t *testing.T) {
	repo := NewRepo(db)
	id := uuid.New().String()
	t.Run("create card", func(t *testing.T) {
		deck := models.Deck{
			UUID:     id,
			Shuffled: false,
			Cards:    "AS,2S,3S,4S,5S,6S,7S,8S,9S,10S,JS,QS,KS,AD,2D,3D,4D,5D,6D,7D,8D,9D,10D,JD,QD,KD,AC,2C,3C,4C,5C,6C,7C,8C,9C,10C,JC,QC,KC,AH,2H,3H,4H,5H,6H,7H,8H,9H,10H,JH,QH,KH",
		}

		err := repo.CreateCard(deck)
		assert.NoError(t, err, expectedStr(nil, err))
	})
	t.Run("get card", func(t *testing.T) {
		deck := models.Deck{
			UUID:     id,
			Shuffled: false,
			Cards:    "AS,2S,3S,4S,5S,6S,7S,8S,9S,10S,JS,QS,KS,AD,2D,3D,4D,5D,6D,7D,8D,9D,10D,JD,QD,KD,AC,2C,3C,4C,5C,6C,7C,8C,9C,10C,JC,QC,KC,AH,2H,3H,4H,5H,6H,7H,8H,9H,10H,JH,QH,KH",
		}

		dbDeck, err := repo.GetDeck(id)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, deck.UUID, dbDeck.UUID, expectedStr(deck.UUID, dbDeck.UUID))
		assert.Equal(t, deck.Shuffled, dbDeck.Shuffled, expectedStr(deck.Shuffled, dbDeck.Shuffled))
		assert.Equal(t, deck.Cards, dbDeck.Cards, expectedStr(deck.Cards, dbDeck.Cards))

	})

	t.Run("update card", func(t *testing.T) {
		deck := models.Deck{
			UUID:     id,
			Shuffled: false,
			Cards:    "AS,2S,3S,4S",
		}

		err := repo.UpdateDeck(deck)

		dbDeck, _ := repo.GetDeck(id)
		assert.NoError(t, err, expectedStr(nil, err))
		assert.Equal(t, dbDeck.UUID, dbDeck.UUID, expectedStr(dbDeck.UUID, dbDeck.UUID))
		assert.Equal(t, dbDeck.Cards, dbDeck.Cards, expectedStr(dbDeck.Cards, dbDeck.Cards))
		assert.Equal(t, dbDeck.Shuffled, dbDeck.Shuffled, expectedStr(dbDeck.Shuffled, dbDeck.Shuffled))
	})
}

func expectedStr(expected, got interface{}) string {
	return fmt.Sprintf("Expected %v but got %v", expected, got)
}

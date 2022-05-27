package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mustafasegf/decks/internal/logger"
	"github.com/mustafasegf/decks/pkg/api"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	godotenv.Load()
	logger.SetLogger()

	var db *pgxpool.Pool

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		zap.L().Fatal("DB Connection", zap.Error(err))
	}

	defer func() {
		db.Close()
	}()

	s := api.MakeServer(db)
	s.RunServer()
}

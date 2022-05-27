package decks

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mustafasegf/decks/pkg/models"
)

type Repo interface {
	CreateCard(deck models.Deck) (err error)
	GetDeck(uuid string) (deck models.Deck, err error)
	UpdateDeck(deck models.Deck) (err error)
}

type repo struct {
	DB *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *repo {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateCard(deck models.Deck) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert("decks").Columns("uuid", "shuffled", "cards").Values(deck.UUID, deck.Shuffled, deck.Cards).ToSql()
	if err != nil {
		err = fmt.Errorf("could not generate query: %w", err)
	}

	ctx := context.Background()

	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("could not begin transaction: %w", err)
		return
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		err = fmt.Errorf("could not insert card: %w", err)
		return
	}

	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("could not commit transaction: %w", err)
	}
	return
}

func (r *repo) GetDeck(uuid string) (deck models.Deck, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("uuid", "shuffled", "cards").From("decks").Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		err = fmt.Errorf("could not generate query: %w", err)
		return
	}

	ctx := context.Background()
	err = pgxscan.Get(ctx, r.DB, &deck, sql, args...)
	if err != nil {
		err = fmt.Errorf("could not get deck: %w", err)
	}

	return
}

func (r *repo) UpdateDeck(deck models.Deck) (err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("decks").Set("cards", deck.Cards).Set("shuffled", deck.Shuffled).Where(sq.Eq{"uuid": deck.UUID}).ToSql()
	if err != nil {
		err = fmt.Errorf("could not generate query: %w", err)
		return
	}

	ctx := context.Background()

	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		err = fmt.Errorf("could not begin transaction: %w", err)
		return
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		err = fmt.Errorf("could not insert card: %w", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("could not commit transaction: %w", err)
	}
	return
}

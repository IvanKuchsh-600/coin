package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"currency/internal/entities"
	"currency/internal/usecases"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, connStr string) (*Storage, error) {
	fmt.Printf("Connecting to %s\n", connStr)
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database")
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(entities.ErrInternalServer, "Could not connect to pool")
	}
	return &Storage{db: pool}, nil
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	query := `INSERT INTO coins (title, price, created_at) VALUES ($1, $2, $3);`
	for _, coin := range coins {
		_, err := s.db.Exec(ctx, query, coin.Title, coin.Price, coin.CreateTime)
		if err != nil {
			return errors.Wrap(entities.ErrInternalServer, "Coin was not added")
		}
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, titles []string, options ...usecases.Option) ([]entities.Coin, error) {
	fmt.Println(titles)
	opts := &usecases.Options{}
	for _, option := range options {
		option(opts)
	}
	var query string
	switch opts.FuncType {
	case usecases.Max:
		query = `SELECT title, MAX(price), created_at FROM coins WHERE title = $1 GROUP BY title;`
	case usecases.Min:
		query = `SELECT title, MIN(price), created_at FROM coins WHERE title = $1 GROUP BY title;`
	case usecases.Avg:
		query = `SELECT title, AVG(price), created_at FROM coins WHERE title = $1 GROUP BY title;`
	default:
		query = `SELECT title, price, created_at FROM coins WHERE title = $1 ORDER BY created_at DESC LIMIT 1;`
	}

	var coin entities.Coin
	var coins []entities.Coin
	for _, t := range titles {
		err := s.db.QueryRow(ctx, query, t).Scan(&coin.Title, &coin.Price, &coin.CreateTime)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.Wrap(entities.ErrInvalidParams, fmt.Sprintf("Unable to get coin: %s", coin.Title))
			}
			return nil, errors.Wrap(entities.ErrInternalServer, fmt.Sprintf("Unable to get coin: %s", coin.Title))
		}
		coins = append(coins, coin)
	}

	return coins, nil
}

func (s *Storage) GetTitles(ctx context.Context) ([]string, error) {
	var titles []string
	query := `SELECT DISTINCT title FROM coins;`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get titles")
	}

	defer rows.Close()

	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unable to get title: %s", title))
		}

		titles = append(titles, title)
	}

	return titles, nil
}

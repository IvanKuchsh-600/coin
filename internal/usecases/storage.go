package usecases

import (
	"context"
	"currency/internal/entities"
)

type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error
	Get(ctx context.Context, titles []string, opt ...Option) ([]entities.Coin, error)
	GetTitles(ctx context.Context) ([]string, error)
}

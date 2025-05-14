package public

import (
	"context"

	"currency/internal/entities"
)

type Service interface {
	GetLastPrice(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMinPrice(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMaxPrice(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvgPrice(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetCoinsFromAPI(ctx context.Context, titles ...string) ([]entities.Coin, error)
}

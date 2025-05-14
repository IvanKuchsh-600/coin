package usecases

import (
	"context"
	"currency/internal/entities"
)

type Client interface {
	GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
}

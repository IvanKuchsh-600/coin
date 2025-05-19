package usecases

import (
	"context"
	"currency/internal/entities"
)

//go:generate mockgen -source=client.go -destination=./mocks/client_mock.go -package=mock
type Client interface {
	GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
}

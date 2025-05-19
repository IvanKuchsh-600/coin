package usecases

import (
	"context"
	"fmt"

	"currency/internal/entities"

	"github.com/pkg/errors"
)

type Service struct {
	storage Storage
	client  Client
}

func NewService(storage Storage, client Client) (*Service, error) {
	if storage == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "storage is nil")
	}
	if client == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}

type AggFunc int

const (
	_ AggFunc = iota
	Max
	Min
	Avg
)

type Options struct {
	FuncType AggFunc
}

type Option func(opts *Options)

func (af AggFunc) String() string {
	return [...]string{"", "MAX", "MIN", "AVG"}[af]
}

func WithMaxFunc() Option {
	return func(opts *Options) {
		opts.FuncType = Max
	}
}

func WithMinFunc() Option {
	return func(opts *Options) {
		opts.FuncType = Min
	}
}

func WithAvgFunc() Option {
	return func(opts *Options) {
		opts.FuncType = Avg
	}
}

func (s *Service) GetLastPrice(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles)
	if errors.Is(err, entities.ErrInvalidParams) {
		return nil, errors.Wrap(entities.ErrInvalidParams, "incorrect parameters")
	}
	if err != nil {
		return nil, errors.Wrap(entities.ErrGetFunc, "GetLastPrice")
	}

	return coins, nil
}

func (s *Service) GetMaxPrice(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithMaxFunc())
	if errors.Is(err, entities.ErrInvalidParams) {
		return nil, errors.Wrap(entities.ErrInvalidParams, "incorrect parameters")
	}

	if err != nil {
		return nil, errors.Wrap(entities.ErrGetFunc, "GetMaxPrice")
	}

	return coins, nil
}

func (s *Service) GetMinPrice(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithMinFunc())
	if errors.Is(err, entities.ErrInvalidParams) {
		return nil, errors.Wrap(entities.ErrInvalidParams, "incorrect parameters")
	}

	if err != nil {
		return nil, errors.Wrap(entities.ErrGetFunc, "GetMinPrice")
	}

	return coins, nil
}

func (s *Service) GetAvgPrice(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithAvgFunc())
	if errors.Is(err, entities.ErrInvalidParams) {
		return nil, errors.Wrap(entities.ErrInvalidParams, "incorrect parameters")
	}

	if err != nil {

		return nil, errors.Wrap(entities.ErrGetFunc, "GetAvgPrice")
	}

	return coins, nil
}

func (s *Service) GetCoinsFromAPI(ctx context.Context, titles ...string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		ts, err := s.storage.GetTitles(ctx)
		if err != nil {
			return nil, errors.Wrap(entities.ErrGetFunc, "GetCoinsFromAPI")
		}
		titles = ts
	}

	coins, err := s.client.GetCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrGetFunc, "GetCoinsFromAPI")
	}

	err = s.storage.Store(ctx, coins)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(entities.ErrGetFunc, "GetCoinsFromAPI")
	}

	return coins, nil
}

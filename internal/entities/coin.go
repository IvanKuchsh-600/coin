package entities

import (
	"time"

	"github.com/pkg/errors"
)

type Coin struct {
	Title      string
	Price      float64
	CreateTime time.Time
}

func NewCoin(title string, price float64, created time.Time) (*Coin, error) {
	if title == "" {
		return nil, errors.Wrap(ErrInvalidParams, "Title is empty")
	}
	if price < 0 {
		return nil, errors.Wrap(ErrInvalidParams, "Price negative")
	}
	return &Coin{Title: title, Price: price, CreateTime: created}, nil
}

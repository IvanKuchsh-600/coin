package coindesk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"currency/internal/entities"

	"github.com/pkg/errors"
)

type Client struct {
	client http.Client
	url    string
}

func NewClient(url string) (*Client, error) {
	cl := http.Client{}
	return &Client{client: cl, url: url}, nil
}

func (c *Client) GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	fsymsParams := strings.Join(titles, ",")
	url := fmt.Sprintf("%s=%s", c.url, fsymsParams)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't form a request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Couldn't get %s", fsymsParams))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status Error: %s\n", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't count the response")
	}

	var priceData map[string]map[string]float64

	err = json.Unmarshal(bodyBytes, &priceData)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, fmt.Sprintf("titles: %s", fsymsParams))
	}

	var coins []entities.Coin
	for coin, prices := range priceData {
		c, err := entities.NewCoin(coin, prices["RUB"], time.Now())
		if err != nil {
			return nil, err
		}
		coins = append(coins, *c)
	}

	return coins, nil
}

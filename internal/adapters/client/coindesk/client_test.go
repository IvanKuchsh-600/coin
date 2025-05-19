package coindesk_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"currency/internal/adapters/client/coindesk"
	"currency/internal/entities"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetCoins(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/?fsyms=BTC,ETH", r.URL.String())
			assert.Equal(t, http.MethodGet, r.Method)

			response := `{
				"BTC": {"RUB": 8398290.1},
				"ETH": {"RUB": 196888.49}
        	}`

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(response))
		}))
		defer testServer.Close()

		client, err := coindesk.NewClient(testServer.URL + "?fsyms")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		coins, err := client.GetCoins(context.Background(), []string{"BTC", "ETH"})
		if err != nil {
			t.Errorf("Error getting coins: %v", err)
			return
		}

		// Проверяем результаты
		expected := []entities.Coin{
			{Title: "BTC", Price: 8398290.1, CreateTime: time.Now()},
			{Title: "ETH", Price: 196888.49, CreateTime: time.Now()},
		}

		assert.Len(t, coins, 2)
		assert.Equal(t, expected[0].Title, coins[0].Title)
		assert.Equal(t, expected[0].Price, coins[0].Price)
		assert.Equal(t, expected[1].Title, coins[1].Title)
		assert.Equal(t, expected[1].Price, coins[1].Price)
	})

	t.Run("empty response", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}))
		defer testServer.Close()

		client, err := coindesk.NewClient(testServer.URL + "?fsyms")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}

		coins, err := client.GetCoins(context.Background(), []string{"BTC", "ETH"})
		if err != nil {
			t.Errorf("Error getting coins: %v", err)
			return
		}
		assert.Empty(t, coins)
	})

	t.Run("invalid status code", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer testServer.Close()

		client, err := coindesk.NewClient(testServer.URL + "?fsyms")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}

		_, err = client.GetCoins(context.Background(), []string{"BTC", "ETH"})
		require.Error(t, err)

		assert.Contains(t, err.Error(), "Status Error: 500 Internal Server Error")
	})

	t.Run("invalid json", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer testServer.Close()

		client, err := coindesk.NewClient(testServer.URL + "?fsyms")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		_, err = client.GetCoins(context.Background(), []string{"BTCCC", "ETGFH"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid Params")
	})

	t.Run("request error", func(t *testing.T) {
		client, err := coindesk.NewClient("http://invalid-url")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		_, err = client.GetCoins(context.Background(), []string{"BTC"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Couldn't get BTC")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Сервер, который долго отвечает
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer testServer.Close()

		client, err := coindesk.NewClient(testServer.URL + "?fsyms")
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err = client.GetCoins(ctx, []string{"BTC"})
		require.Error(t, err)
		fmt.Println(err.Error())
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}

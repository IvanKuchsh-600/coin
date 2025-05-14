package public

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"currency/internal/entities"
	"currency/pkg/dto"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Coin API
//	@version		1.0
//	@description	This is a sample server for Coin API.
//	@license.name	Apache 2.0
//  @license.url http://www.apache.org/licenses/LICENSE-2.0.html

//	@host			localhost:8080
//	@BasePath		/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

type Server struct {
	port    string
	r       *chi.Mux
	service Service
}

func NewServer(service Service, port string) (*Server, error) {
	if service == nil || service == Service(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParams, "service is nil")
	}

	r := chi.NewRouter()
	return &Server{port: port, r: r, service: service}, nil
}

func (s *Server) Run() error {
	s.r.Get("/v1/get_current_rate", s.GetLastPriceHandler)
	s.r.Get("/v1/get_max_rate", s.GetMaxPriceHandler)
	s.r.Get("/v1/get_min_rate", s.GetMinPriceHandler)
	s.r.Get("/v1/get_avg_rate", s.GetAvgPriceHandler)

	s.r.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))
	s.r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.r)
	if err != nil {
		return errors.Wrap(entities.ErrInternalServer, err.Error())
	}

	return nil
}

// GetLastPriceHandler godoc
//
//	@Summary		Get current rate
//	@Description	Get the current rate of specified coins
//	@Tags			coins
//	@Accept			json
//	@Produce		json
//	@Param			fsyms	query		string	true	"Comma-separated list of cryptocurrencies"
//	@Success		200		{object}		dto.CoinsDTO "List of cryptocurrencies"
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/v1/get_current_rate [get]
func (s *Server) GetLastPriceHandler(rw http.ResponseWriter, req *http.Request) {
	titles := strings.Split(req.URL.Query().Get("fsyms"), ",")
	ctx := req.Context()

	coins, err := s.service.GetLastPrice(ctx, titles) // Пытаемся взять из БД
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParams) { // В БД не нашли таких titles
			cs, err := s.service.GetCoinsFromAPI(ctx, titles...) // Берем из API
			if err != nil {
				fmt.Println(err)
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			coins = cs
		}
		if !errors.Is(err, entities.ErrInvalidParams) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var coinsDTO dto.CoinsDTO
	for _, coin := range coins {
		coinsDTO = append(coinsDTO, dto.CoinDTO{
			Title:      coin.Title,
			Price:      math.Round(coin.Price*100) / 100,
			CreateTime: coin.CreateTime.Format("2006-02-02"),
		})
	}

	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(coinsDTO); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetMaxPriceHandler godoc
//
//	@Summary		Get max rate
//	@Description	Get the max rate of specified coins
//	@Tags			coins
//	@Accept			json
//	@Produce		json
//	@Param			fsyms	query		string	true	"Comma-separated list of cryptocurrencies"
//	@Success		200		{array}		dto.CoinDTO
//	@Failure		400		{object}	map[string]interface{}	"Invalid input"
//	@Failure		404		{object}	map[string]interface{}	"No coins found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/get_max_rate [get]
func (s *Server) GetMaxPriceHandler(rw http.ResponseWriter, req *http.Request) {
	titles := strings.Split(req.URL.Query().Get("fsyms"), ",")
	ctx := req.Context()

	coins, err := s.service.GetLastPrice(ctx, titles)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParams) {
			cs, err := s.service.GetCoinsFromAPI(ctx, titles...)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			}
			coins = cs
		}
		if !errors.Is(err, entities.ErrInvalidParams) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var coinsDTO dto.CoinsDTO
	for _, coin := range coins {
		coinsDTO = append(coinsDTO, dto.CoinDTO{
			Title:      coin.Title,
			Price:      coin.Price,
			CreateTime: coin.CreateTime.Format("2006-02-02"),
		})
	}

	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(coinsDTO); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//rw.WriteHeader(http.StatusOK)
}

// GetMinPriceHandler godoc
//
//	@Summary		Get min rate
//	@Description	Get the min rate of specified coins
//	@Tags			coins
//	@Accept			json
//	@Produce		json
//	@Param			fsyms	query		string	true	"Comma-separated list of cryptocurrencies"
//	@Success		200		{array}		dto.CoinDTO
//	@Failure		400		{object}	map[string]interface{}	"Invalid input"
//	@Failure		404		{object}	map[string]interface{}	"No coins found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/get_min_rate [get]
func (s *Server) GetMinPriceHandler(rw http.ResponseWriter, req *http.Request) {
	titles := strings.Split(req.URL.Query().Get("fsyms"), ",")
	ctx := req.Context()

	coins, err := s.service.GetLastPrice(ctx, titles)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParams) {
			cs, err := s.service.GetCoinsFromAPI(ctx, titles...)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			}
			coins = cs
		}
		if !errors.Is(err, entities.ErrInvalidParams) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var coinsDTO dto.CoinsDTO
	for _, coin := range coins {
		coinsDTO = append(coinsDTO, dto.CoinDTO{
			Title:      coin.Title,
			Price:      coin.Price,
			CreateTime: coin.CreateTime.Format("2006-02-02"),
		})
	}

	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(coinsDTO); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAvgPriceHandler godoc
//
//	@Summary		Get avg rate
//	@Description	Get the avg rate of specified coins
//	@Tags			coins
//	@Accept			json
//	@Produce		json
//	@Param			fsyms	query		string	true	"Comma-separated list of cryptocurrencies"
//	@Success		200		{array}		dto.CoinDTO
//	@Failure		400		{object}	map[string]interface{}	"Invalid input"
//	@Failure		404		{object}	map[string]interface{}	"No coins found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/get_avg_rate [get]
func (s *Server) GetAvgPriceHandler(rw http.ResponseWriter, req *http.Request) {
	titles := strings.Split(req.URL.Query().Get("fsyms"), ",")
	ctx := req.Context()

	coins, err := s.service.GetLastPrice(ctx, titles)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParams) {
			cs, err := s.service.GetCoinsFromAPI(ctx, titles...)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			}
			coins = cs
		}
		if !errors.Is(err, entities.ErrInvalidParams) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var coinsDTO dto.CoinsDTO
	for _, coin := range coins {
		coinsDTO = append(coinsDTO, dto.CoinDTO{
			Title:      coin.Title,
			Price:      coin.Price,
			CreateTime: coin.CreateTime.Format("2006-02-02"),
		})
	}

	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(coinsDTO); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

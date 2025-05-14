package dto

type CoinDTO struct {
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	CreateTime string  `json:"create_time"`
}

type CoinsDTO []CoinDTO

package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Pair struct {
	Quote string `json:"quote"`
}

type LatestResponse struct {
	Rate        decimal.Decimal `json:"rate"`
	LastUpdated time.Time       `json:"last_updated"`
}

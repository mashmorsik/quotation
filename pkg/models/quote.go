package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Quote struct {
	ID             uuid.UUID       `json:"id"`
	BaseCurrency   string          `json:"base_currency"`
	TargetCurrency string          `json:"target_currency"`
	Timestamp      time.Time       `json:"timestamp"`
	Rate           decimal.Decimal `json:"rate"`
}

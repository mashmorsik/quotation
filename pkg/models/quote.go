package models

import (
	"github.com/docker/distribution/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Quote struct {
	ID             uuid.UUID
	BaseCurrency   string
	TargetCurrency string
	Timestamp      time.Time
	Rate           decimal.Decimal
}

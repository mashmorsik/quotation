package models

import "github.com/shopspring/decimal"

type QuoteAPIResponse struct {
	Base  string                     `json:"base"`
	Date  string                     `json:"date"`
	Rates map[string]decimal.Decimal `json:"rates"`
}

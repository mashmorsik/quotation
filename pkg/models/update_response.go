package models

import "github.com/google/uuid"

type UpdateResponse struct {
	QuoteID uuid.UUID `json:"quoteID"`
}

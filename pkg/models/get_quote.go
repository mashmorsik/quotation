package models

import "github.com/google/uuid"

type GetRequest struct {
	ID uuid.UUID `json:"quote_id"`
}

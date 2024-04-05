package repository

import (
	"github.com/google/uuid"
	"github.com/mashmorsik/quotation/pkg/models"
)

type Repository interface {
	AddQuotePair(from, to string) error
	GetQuotePairs() ([][]string, error)
	AddQuotation(q *models.Quote) error
	GetQuotation(id uuid.UUID) (*models.Quote, error)
	GetLastUpdated(from, to string) (*models.Quote, error)
}

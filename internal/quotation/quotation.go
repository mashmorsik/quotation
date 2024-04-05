package quotation

import (
	"context"
	"github.com/google/uuid"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/infrastructure/quote_api"
	"github.com/mashmorsik/quotation/pkg/models"
	"github.com/mashmorsik/quotation/repository"
	errs "github.com/pkg/errors"
	"time"
)

type Quotation struct {
	Ctx    context.Context
	Repo   repository.Repository
	Config *config.Config
}

func NewQuotation(ctx context.Context, repo repository.Repository, conf *config.Config) *Quotation {
	return &Quotation{Ctx: ctx, Repo: repo, Config: conf}
}

func (q *Quotation) GetQuoteAsync(from, to string) (uuid.UUID, error) {
	quoteID := uuid.New()

	rate, err := quote_api.GetQuote(from, to, q.Config)
	if err != nil {
		return uuid.UUID{}, errs.WithMessagef(err, "failed to get quote for %s/%s", from, to)
	}

	quote := &models.Quote{
		ID:             quoteID,
		BaseCurrency:   from,
		TargetCurrency: to,
		Timestamp:      time.Now(),
		Rate:           rate,
	}

	err = q.Repo.AddQuotePair(quote.BaseCurrency, quote.TargetCurrency)
	if err != nil {
		return uuid.UUID{}, errs.WithMessagef(err, "failed to AddQuotePair, for: %v\n", quote)
	}

	quoteLatest, err := q.Repo.GetLastUpdated(quote.BaseCurrency, quote.TargetCurrency)
	if err != nil {
		return uuid.UUID{}, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v\n", quote)
	}

	now := time.Now().UTC()
	if quoteLatest != nil {
		if quoteLatest.Timestamp.Add(q.Config.ResponseDelay).After(now) {
			return quoteLatest.ID, nil
		}
	}

	if err = q.Repo.AddQuotation(quote); err != nil {
		return uuid.UUID{}, errs.WithMessagef(err, "failed to AddQuotation, for: %v\n", quote)
	}
	return quoteID, nil
}

func (q *Quotation) GetQuotationByID(quoteID uuid.UUID) (*models.Quote, error) {
	quote, err := q.Repo.GetQuotation(quoteID)
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to GetQuotationByID, for: %v", quoteID)
	}

	return quote, nil
}

func (q *Quotation) GetLastUpdated(from, to string) (*models.Quote, error) {
	quotePairs, err := q.Repo.GetQuotePairs()
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", from)
	}

	if len(quotePairs) == 0 {
		quoteID, err := q.GetQuoteAsync(from, to)
		if err != nil {
			return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", from)
		}
		quote, err := q.GetQuotationByID(quoteID)
		if err != nil {
			return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", quoteID)
		}
		return quote, nil
	}

	for _, pair := range quotePairs {
		if pair[0] == from && pair[1] == to {
			quote, err := q.Repo.GetLastUpdated(from, to)
			if err != nil {
				return nil, errs.WithMessagef(err, "failed to GetLastUpdated for %s/%s", from, to)
			}
			return quote, nil
		} else {
			quoteID, err := q.GetQuoteAsync(from, to)
			if err != nil {
				return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", from)
			}
			quote, err := q.GetQuotationByID(quoteID)
			if err != nil {
				return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", quoteID)
			}
			return quote, nil
		}
	}

	return nil, errs.WithMessagef(err, "failed to GetLastUpdated, for: %v", from)
}

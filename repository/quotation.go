package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/infrastructure/data"
	"github.com/mashmorsik/quotation/pkg/models"
	errs "github.com/pkg/errors"
	"time"
)

type QuoteRepo struct {
	Ctx  context.Context
	data *data.Data
}

func NewQuoteRepo(ctx context.Context, data *data.Data) *QuoteRepo {
	return &QuoteRepo{Ctx: ctx, data: data}
}

func (qr *QuoteRepo) AddQuotePair(from, to string) error {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	qPair, err := qr.GetQuotePair(from, to)
	if err != nil {
		return err
	}

	if qPair != nil {
		return nil
	}

	res, err := qr.data.Master().ExecContext(ctx, `
		INSERT INTO quote_pair (base_currency, target_currency)
		VALUES ($1, $2)`, from, to)
	if err != nil {
		return errs.WithMessagef(err, "failed to add quote pair to database, "+
			"from: %s, to: %s\n", from, to)
	}
	ra, _ := res.RowsAffected()
	logger.Infof("rows affected: %v", ra)

	return nil
}

func (qr *QuoteRepo) GetQuotePair(from, to string) (*int32, error) {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	var quotePairId int32

	err := qr.data.Master().QueryRowContext(ctx,
		`SELECT id 
		FROM quote_pair
		WHERE base_currency = $1 AND target_currency = $2`, from, to).Scan(&quotePairId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errs.WithMessagef(err, "failed to exec query: GetQuotePair")
	}

	return &quotePairId, nil
}

func (qr *QuoteRepo) GetQuotePairs() ([][]string, error) {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	var quotePairs [][]string

	query := `
		SELECT base_currency, target_currency
		FROM quote_pair`

	rows, err := qr.data.Master().QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.WithMessage(err, "no quote pairs found")
		} else {
			return nil, errs.WithMessagef(err, "failed to exec query: %s", query)
		}
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Errf("failed to close rows: %s", err.Error())
			return
		}
	}(rows)

	for rows.Next() {
		var quotePair []string
		var baseCurrency string
		var targetCurrency string

		if err = rows.Scan(&baseCurrency, &targetCurrency); err != nil {
			return nil, errs.WithMessagef(err, "failed to scan row")
		}
		quotePair = append(quotePair, baseCurrency, targetCurrency)
		quotePairs = append(quotePairs, quotePair)
	}

	return quotePairs, nil
}

func (qr *QuoteRepo) AddQuotation(q *models.Quote) error {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	query := `
		INSERT INTO quotation (id, base_currency, target_currency, rate, time_updated) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := qr.data.Master().ExecContext(ctx, query, q.ID, q.BaseCurrency, q.TargetCurrency, q.Rate, q.Timestamp)
	if err != nil {
		return errs.WithMessagef(err, "failed to add quote for quoteID: %s", q.ID)
	}

	return nil
}

func (qr *QuoteRepo) GetQuotation(id uuid.UUID) (*models.Quote, error) {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	var q models.Quote

	query := `
		SELECT id, base_currency, target_currency, rate, time_updated
		FROM quotation
		WHERE id = $1`

	err := qr.data.Master().QueryRowContext(ctx, query, id).Scan(&q.ID, &q.BaseCurrency, &q.TargetCurrency, &q.Rate, &q.Timestamp)
	if err != nil {
		return nil, errs.WithMessagef(err, "failed to get quote for id: %s", id)
	}

	return &q, nil
}

func (qr *QuoteRepo) GetLastUpdated(from, to string) (*models.Quote, error) {
	ctx, cancel := context.WithTimeout(qr.Ctx, time.Second*5)
	defer cancel()

	var q models.Quote

	err := qr.data.Master().QueryRowContext(ctx, `
		SELECT * 
		FROM quotation
		WHERE base_currency = $1 AND target_currency = $2
		ORDER BY time_updated DESC 
		LIMIT 1`, from, to).
		Scan(&q.ID, &q.BaseCurrency, &q.TargetCurrency, &q.Rate, &q.Timestamp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.WithMessagef(err, "failed to get quote for %s/%s", from, to)
	}

	return &q, nil
}

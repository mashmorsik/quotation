package server

import (
	"errors"
	"github.com/mashmorsik/quotation/pkg/currency"
	"slices"
	"strings"
)

func (s *HTTPServer) validateQuote(quote string) error {
	if quote == "" {
		return errors.New("currency pair is required")
	}
	if len([]rune(quote)) < 6 {
		return errors.New("currency pair is too short")
	}

	if !strings.Contains(quote, "/") {
		return errors.New("currency pair requires separator")
	}

	from, to := currency.SeparateCurrency(quote)

	if from == to {
		return errors.New("currencies are the same")
	}

	if !slices.Contains(s.Config.Quotations, from) || !slices.Contains(s.Config.Quotations, to) {
		return errors.New("currency pair is invalid")
	}

	return nil
}

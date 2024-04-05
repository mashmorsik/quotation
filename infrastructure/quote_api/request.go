package quote

import (
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/pkg/models"
	errs "github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
)

func GetQuote(from, to string) (decimal.Decimal, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return decimal.Zero, errs.WithMessage(err, "failed to load config")
	}

	reqStr := fmt.Sprintf(conf.QuoteAPI, from, to)
	req, err := http.NewRequest("GET", reqStr, nil)
	if err != nil {
		return decimal.Zero, errs.WithMessage(err, "failed to create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, errs.WithMessagef(err, "failed to do request: %v", req)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Errf("failed to close response body: %v", err)
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return decimal.Zero, errs.WithMessagef(err, "invalid response status: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return decimal.Zero, errs.WithMessage(err, "failed to read response body")
	}

	var response models.Response
	if err = json.Unmarshal(body, &response); err != nil {
		return decimal.Zero, errs.WithMessagef(err, "failed to unmarshal response, body: %v", body)
	}

	rate, found := response.Rates[to]
	if !found {
		return decimal.Zero, errs.WithMessagef(err, "failed to find EUR rate")
	}

}

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/internal/quotation"
	"github.com/mashmorsik/quotation/pkg/models"
	mock_repository "github.com/mashmorsik/quotation/test/testdata/mock_repo"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHTTPServer_UpdateQuote(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	latestID := uuid.New()

	from := "EUR"
	to := "MXN"

	latestQuote := &models.Quote{
		ID:             latestID,
		BaseCurrency:   "EUR",
		TargetCurrency: "MXN",
		Timestamp:      time.Now().UTC(),
		Rate:           decimal.NewFromFloat(0.876),
	}

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().AddQuotePair(from, to).Return(nil)
	mockRepo.EXPECT().GetLastUpdated(from, to).Return(latestQuote, nil)

	conf := &config.Config{
		ResponseDelay: 5 * time.Second,
		Quotations:    []string{"EUR", "MXN", "USD"},
		QuoteAPI: struct {
			URL string `yaml:"url"`
		}{
			URL: "https://api.frankfurter.app/latest?from=%s&to=%s"},
	}
	q := &quotation.Quotation{
		Ctx:    context.Background(),
		Repo:   mockRepo,
		Config: conf,
	}
	srv := NewServer(conf, *q)
	testServer := httptest.NewServer(http.HandlerFunc(srv.UpdateQuote))
	defer testServer.Close()

	reqPair := &models.Pair{Quote: "EUR/MXN"}
	jsonStr, err := json.Marshal(reqPair)
	if err != nil {
		t.Fatalf("failed to marshal reqPair: %v", err)
	}
	resp, err := http.Post(testServer.URL+"/update", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("Error getting quote: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %v", resp.Header.Get("Content-Type"))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	gotID := &models.UpdateResponse{}
	err = json.Unmarshal(respBody, gotID)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	if gotID.QuoteID != latestQuote.ID {
		t.Errorf("Wanted: %+v, got: %+v", latestQuote.ID, gotID)
	}
}

func TestHTTPServer_GetQuote(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()

	quote := &models.Quote{
		ID:             id,
		BaseCurrency:   "EUR",
		TargetCurrency: "USD",
		Timestamp:      time.Time{},
		Rate:           decimal.NewFromFloat(0.876),
	}

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetQuotation(id).Return(quote, nil)

	conf := &config.Config{
		ResponseDelay: 2 * time.Second,
	}
	q := &quotation.Quotation{
		Ctx:    context.Background(),
		Repo:   mockRepo,
		Config: conf,
	}
	srv := NewServer(conf, *q)
	testServer := httptest.NewServer(http.HandlerFunc(srv.GetQuote))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/get?quoteID=" + id.String())
	if err != nil {
		t.Errorf("Error getting quote: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %v", resp.Header.Get("Content-Type"))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	gotQuote := &models.Quote{}
	err = json.Unmarshal(respBody, gotQuote)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	if !reflect.DeepEqual(gotQuote, quote) {
		t.Errorf("Unexpected quote: %v", gotQuote)
	}
}

func TestHTTPServer_GetLatestQuote_quote_not_in_database(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	latestID := uuid.New()

	from := "USD"
	to := "MXN"
	timestamp := time.Now().UTC()

	latestQuote := &models.Quote{
		ID:             latestID,
		BaseCurrency:   "USD",
		TargetCurrency: "MXN",
		Timestamp:      timestamp,
		Rate:           decimal.NewFromFloat(17.03),
	}

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetQuotePairs().Return([][]string{}, nil)
	mockRepo.EXPECT().AddQuotePair(from, to).Return(nil)
	mockRepo.EXPECT().GetLastUpdated(from, to).Return(latestQuote, nil)
	mockRepo.EXPECT().GetQuotation(latestID).Return(latestQuote, nil)

	conf := &config.Config{
		ResponseDelay: 5 * time.Second,
		Quotations:    []string{"EUR", "MXN", "USD"},
		QuoteAPI: struct {
			URL string `yaml:"url"`
		}{
			URL: "https://api.frankfurter.app/latest?from=%s&to=%s"},
	}
	q := &quotation.Quotation{
		Ctx:    context.Background(),
		Repo:   mockRepo,
		Config: conf,
	}
	srv := NewServer(conf, *q)
	testServer := httptest.NewServer(http.HandlerFunc(srv.GetLatestQuote))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/latest?quote=USD/MXN")
	if err != nil {
		t.Errorf("Error getting quote: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %v", resp.Header.Get("Content-Type"))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	got := &models.LatestResponse{}
	err = json.Unmarshal(respBody, got)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	if !reflect.DeepEqual(got.Rate, latestQuote.Rate) || !reflect.DeepEqual(got.LastUpdated, latestQuote.Timestamp) {
		t.Errorf("Wanted rate: %+v, got rate: %+v, wanted timestamp: %+v, got timestamp: %+v",
			latestQuote.Rate, got.Rate, latestQuote.Timestamp, got.LastUpdated)
	}
}

func TestHTTPServer_GetLatestQuote_quote_in_database(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	latestID := uuid.New()

	from := "USD"
	to := "MXN"
	timestamp := time.Now().UTC()

	latestQuote := &models.Quote{
		ID:             latestID,
		BaseCurrency:   "USD",
		TargetCurrency: "MXN",
		Timestamp:      timestamp,
		Rate:           decimal.NewFromFloat(17.03),
	}

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetQuotePairs().Return([][]string{{"USD", "MXN"}}, nil)
	mockRepo.EXPECT().GetLastUpdated(from, to).Return(latestQuote, nil)

	conf := &config.Config{
		ResponseDelay: 5 * time.Second,
		Quotations:    []string{"EUR", "MXN", "USD"},
		QuoteAPI: struct {
			URL string `yaml:"url"`
		}{
			URL: "https://api.frankfurter.app/latest?from=%s&to=%s"},
	}
	q := &quotation.Quotation{
		Ctx:    context.Background(),
		Repo:   mockRepo,
		Config: conf,
	}
	srv := NewServer(conf, *q)
	testServer := httptest.NewServer(http.HandlerFunc(srv.GetLatestQuote))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/latest?quote=USD/MXN")
	if err != nil {
		t.Errorf("Error getting quote: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %v", resp.Header.Get("Content-Type"))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	got := &models.LatestResponse{}
	err = json.Unmarshal(respBody, got)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	if !reflect.DeepEqual(got.Rate, latestQuote.Rate) || !reflect.DeepEqual(got.LastUpdated, latestQuote.Timestamp) {
		t.Errorf("Wanted rate: %+v, got rate: %+v, wanted timestamp: %+v, got timestamp: %+v",
			latestQuote.Rate, got.Rate, latestQuote.Timestamp, got.LastUpdated)
	}
}

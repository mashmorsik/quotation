package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/internal/quotation"
	"github.com/mashmorsik/quotation/pkg/currency"
	"github.com/mashmorsik/quotation/pkg/middleware"
	"github.com/mashmorsik/quotation/pkg/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type HTTPServer struct {
	Config *config.Config
	Quote  quotation.Quotation
}

func NewServer(conf *config.Config, quote quotation.Quotation) *HTTPServer {
	return &HTTPServer{Config: conf, Quote: quote}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.HandleFunc("/update", middleware.LoggingMiddleware(s.updateQuote)).Methods(http.MethodPost)
	router.HandleFunc("/get", middleware.LoggingMiddleware(s.getQuote)).Methods(http.MethodGet)
	router.HandleFunc("/latest", middleware.LoggingMiddleware(s.getLatestQuote)).Methods(http.MethodGet)

	logger.Infof("HTTPServer is listening on port: %s\n", s.Config.Server.Port)

	err := http.ListenAndServe(s.Config.Server.Port, router)
	if err != nil {
		return errors.WithMessagef(err, "server can't ListenAndServe http requests")
	}

	return nil
}

func (s *HTTPServer) updateQuote(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var reqBody models.Pair
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	if err = s.validateQuote(reqBody.Quote); err != nil {
		logger.Errf("invalid quotePair: %s", reqBody.Quote)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, to := currency.SeparateCurrency(reqBody.Quote)

	quoteID, err := s.Quote.GetQuoteAsync(from, to)
	if err != nil {
		logger.Errf("fail to GetQuoteAsync, for %s/%s", from, to)
		http.Error(w, "fail to GetQuoteAsync", http.StatusInternalServerError)
		return
	}

	response := models.UpdateResponse{QuoteID: quoteID}

	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.Errf("failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		logger.Errf("failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) getQuote(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var reqBody models.GetRequest
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	quoteID := reqBody.ID

	quote, err := s.Quote.GetQuotationByID(quoteID)
	if err != nil {
		http.Error(w, "Failed to get quote", http.StatusNotFound)
	}

	jsonData, err := json.Marshal(quote)
	if err != nil {
		logger.Errf("failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		logger.Errf("failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) getLatestQuote(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var reqBody models.Pair
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	quoteArr := strings.Split(reqBody.Quote, "/")

	quote, err := s.Quote.GetLastUpdated(quoteArr[0], quoteArr[1])
	if err != nil {
		errStr := fmt.Sprintf("Fail to get last updated quote: %v", err)
		http.Error(w, errStr, http.StatusNotFound)
	}

	latestResponse := &models.LatestResponse{
		Rate:        quote.Rate,
		LastUpdated: quote.Timestamp,
	}

	jsonData, err := json.Marshal(latestResponse)
	if err != nil {
		logger.Errf("failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		logger.Errf("failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

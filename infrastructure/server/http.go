package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/internal/quotation"
	"github.com/mashmorsik/quotation/pkg/currency"
	mw "github.com/mashmorsik/quotation/pkg/middleware"
	"github.com/mashmorsik/quotation/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
)

type HTTPServer struct {
	Config *config.Config
	Quote  quotation.Quotation
}

func NewServer(conf *config.Config, quote quotation.Quotation) *HTTPServer {
	return &HTTPServer{Config: conf, Quote: quote}
}

func (s *HTTPServer) StartServer(ctx context.Context) error {

	router := mux.NewRouter()

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	sh := middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/swagger",
		SpecURL: "swagger.yaml",
	}, nil)
	router.Handle("/swagger", sh)

	router.HandleFunc("/update", s.UpdateQuote).Methods(http.MethodPost)
	router.HandleFunc("/get", s.GetQuote).Methods(http.MethodGet)
	router.HandleFunc("/latest", s.GetLatestQuote).Methods(http.MethodGet)

	logger.Infof("HTTPServer is listening on port: %s\n", s.Config.Server.Port)

	router.Use(mw.LoggingMiddleware)

	httpServer := &http.Server{
		Addr:    s.Config.Server.Port,
		Handler: cors.AllowAll().Handler(router),
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		return errors.WithMessagef(err, "exit reason: %s \n", err)
	}

	return nil
}

func (s *HTTPServer) UpdateQuote(w http.ResponseWriter, r *http.Request) {
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

func (s *HTTPServer) GetQuote(w http.ResponseWriter, r *http.Request) {
	quoteID := r.URL.Query().Get("quoteID")
	u, err := uuid.Parse(quoteID)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}

	quote, err := s.Quote.GetQuotationByID(u)
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

func (s *HTTPServer) GetLatestQuote(w http.ResponseWriter, r *http.Request) {
	qPair := r.URL.Query().Get("quote")

	if err := s.validateQuote(qPair); err != nil {
		logger.Errf("invalid quotePair: %s", qPair)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, to := currency.SeparateCurrency(qPair)

	quote, err := s.Quote.GetLastUpdated(from, to)
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

package server

import (
	"github.com/mashmorsik/quotation/config"
	"testing"
)

func TestHTTPServer_validateQuote_errors(t *testing.T) {
	tests := []struct {
		name    string
		quote   string
		wantErr string
	}{
		{
			name:    "empty_quote",
			quote:   "",
			wantErr: "currency pair is required",
		},
		{
			name:    "pair_is_too_short",
			quote:   "EUR/U",
			wantErr: "currency pair is too short",
		},
		{
			name:    "pair_without_separator",
			quote:   "EURUSD",
			wantErr: "currency pair requires separator",
		},
		{
			name:    "invalid_currencies",
			quote:   "GBR/EUR",
			wantErr: "currency pair is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HTTPServer{
				Config: &config.Config{
					Quotations: []string{"EUR", "USD"}},
			}
			if err := s.validateQuote(tt.quote); (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("validateQuote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

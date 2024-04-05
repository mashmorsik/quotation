package server

import (
	"github.com/mashmorsik/quotation/internal/quotation"
	"testing"
)

func TestHTTPServer_validateQuote(t *testing.T) {
	type fields struct {
		Config *config.Config
		Quote  quotation.Quotation
	}
	type args struct {
		quote string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HTTPServer{
				Config: tt.fields.Config,
				Quote:  tt.fields.Quote,
			}
			if err := s.validateQuote(tt.args.quote); (err != nil) != tt.wantErr {
				t.Errorf("validateQuote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

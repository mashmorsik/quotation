package quotation

import (
	"context"
	"github.com/google/uuid"
	"github.com/mashmorsik/quotation/repository"
	"reflect"
	"testing"
)

func TestNewQuotation(t *testing.T) {
	type args struct {
		ctx  context.Context
		repo repository.Repository
		conf *config.Config
	}
	tests := []struct {
		name string
		args args
		want *Quotation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuotation(tt.args.ctx, tt.args.repo, tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuotation_GetLastUpdated(t *testing.T) {
	type fields struct {
		Ctx    context.Context
		Repo   repository.Repository
		Config *config.Config
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Quote
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:    tt.fields.Ctx,
				Repo:   tt.fields.Repo,
				Config: tt.fields.Config,
			}
			got, err := q.GetLastUpdated(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastUpdated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastUpdated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuotation_GetQuotationByID(t *testing.T) {
	type fields struct {
		Ctx    context.Context
		Repo   repository.Repository
		Config *config.Config
	}
	type args struct {
		quoteID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Quote
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:    tt.fields.Ctx,
				Repo:   tt.fields.Repo,
				Config: tt.fields.Config,
			}
			got, err := q.GetQuotationByID(tt.args.quoteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuotationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuotationByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuotation_GetQuoteAsync(t *testing.T) {
	type fields struct {
		Ctx    context.Context
		Repo   repository.Repository
		Config *config.Config
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:    tt.fields.Ctx,
				Repo:   tt.fields.Repo,
				Config: tt.fields.Config,
			}
			got, err := q.GetQuoteAsync(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuoteAsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuoteAsync() got = %v, want %v", got, tt.want)
			}
		})
	}
}

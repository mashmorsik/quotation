package quotation

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/pkg/models"
	mock_repository "github.com/mashmorsik/quotation/test/testdata/mock_repo"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
	"time"
)

func TestQuotation_GetLastUpdated(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetQuotePairs().Return([][]string{{"EUR", "USD"}}, nil)
	mockRepo.EXPECT().GetLastUpdated("EUR", "USD").Return(&models.Quote{
		ID:             uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
		BaseCurrency:   "EUR",
		TargetCurrency: "USD",
		Timestamp:      time.Time{},
		Rate:           decimal.NewFromFloat(1.208),
	}, nil)

	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Quote
		wantErr bool
	}{
		{
			name: "return_quoteID_of_the_last_updated",
			args: args{"EUR", "USD"},
			want: &models.Quote{
				ID:             uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
				BaseCurrency:   "EUR",
				TargetCurrency: "USD",
				Timestamp:      time.Time{},
				Rate:           decimal.NewFromFloat(1.208)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:  context.Background(),
				Repo: mockRepo,
				Config: &config.Config{
					ResponseDelay: 2 * time.Second,
				},
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
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetQuotation(uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84")).Return(&models.Quote{
		ID:             uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
		BaseCurrency:   "EUR",
		TargetCurrency: "USD",
		Timestamp:      time.Time{},
		Rate:           decimal.NewFromFloat(1.208),
	}, nil)

	tests := []struct {
		name    string
		ID      uuid.UUID
		want    *models.Quote
		wantErr bool
	}{
		{
			name: "return_quote_by_quoteID",
			ID:   uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
			want: &models.Quote{
				ID:             uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
				BaseCurrency:   "EUR",
				TargetCurrency: "USD",
				Timestamp:      time.Time{},
				Rate:           decimal.NewFromFloat(1.208)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:    context.Background(),
				Repo:   mockRepo,
				Config: &config.Config{},
			}
			got, err := q.GetQuotationByID(tt.ID)
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
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	quote := &models.Quote{
		ID:             uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
		BaseCurrency:   "EUR",
		TargetCurrency: "USD",
		Timestamp:      time.Now().UTC(),
		Rate:           decimal.NewFromFloat(1.208),
	}

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().AddQuotePair("EUR", "USD").Return(nil)
	mockRepo.EXPECT().GetLastUpdated("EUR", "USD").Return(quote, nil)

	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "return_quoteID_of_the_last_updated",
			args:    args{"EUR", "USD"},
			want:    uuid.MustParse("3f8a26f7-97f8-45a5-bda0-1af96b6b7d84"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quotation{
				Ctx:  context.Background(),
				Repo: mockRepo,
				Config: &config.Config{
					Quotations: []string{"EUR", "USD"},
					QuoteAPI: struct {
						URL string `yaml:"url"`
					}{
						URL: "https://api.frankfurter.app/latest?from=EUR&to=USD",
					},
					ResponseDelay: 5 * time.Second,
				},
			}
			got, err := q.GetQuoteAsync(tt.args.from, tt.args.to)
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

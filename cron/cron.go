package cronSc

import (
	"database/sql"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	"github.com/mashmorsik/quotation/infrastructure/quote_api"
	"github.com/mashmorsik/quotation/pkg/loc"
	"github.com/mashmorsik/quotation/pkg/models"
	"github.com/mashmorsik/quotation/repository"
	errs "github.com/pkg/errors"
	"time"
)

type Scheduler struct {
	sched *gocron.Scheduler
}

type Data struct {
	Repo   repository.Repository
	Config *config.Config
	quotes [][]string
}

func NewScheduler(sched *gocron.Scheduler) *Scheduler {
	return &Scheduler{sched: sched}
}

func NewData(repo repository.Repository, conf *config.Config) *Data {
	return &Data{Repo: repo, Config: conf}
}

func (s *Scheduler) Sc() *gocron.Scheduler {
	return s.sched
}

func StartScheduler() *gocron.Scheduler {
	location, err := loc.GetLoc()
	if err != nil {
		logger.Errf("fail to get loc: %v", err)
		panic(err)
	}
	sched := gocron.NewScheduler(location)
	return sched
}

func (d *Data) RunScheduler() (*gocron.Scheduler, error) {
	sc := NewScheduler(StartScheduler())
	scheduler := sc.Sc()
	defer scheduler.StartAsync()

	_, err := scheduler.Cron(d.Config.Cron.Period).Do(func() {
		quotePairs, err := d.Repo.GetQuotePairs()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Errf("no quote pairs in database: %v", err)
			} else {
				logger.Errf("fail to GetQuotePairs: %v", err)
				return
			}
		}

		for _, pair := range quotePairs {
			rate, err := quote_api.GetQuote(pair[0], pair[1], d.Config)
			if err != nil {
				logger.Errf("fail to GetQuote for pair: %v, err: %s", pair, err)
				return
			}

			quote := &models.Quote{
				ID:             uuid.New(),
				BaseCurrency:   pair[0],
				TargetCurrency: pair[1],
				Timestamp:      time.Now(),
				Rate:           rate,
			}
			err = d.Repo.AddQuotation(quote)
			if err != nil {
				logger.Errf("fail to AddQuotation: %v", err)
				return
			}
		}
	})
	if err != nil {
		return nil, errs.WithMessage(err, "fail to Create CronJob")
	}

	return scheduler, nil
}

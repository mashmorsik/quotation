package loc

import (
	"github.com/mashmorsik/quotation/config"
	errs "github.com/pkg/errors"
	"time"
)

func GetLoc() (*time.Location, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, errs.WithMessage(err, "failed to load config")
	}

	loc, err := time.LoadLocation(conf.Cron.Location)
	if err != nil {
		return nil, errs.WithMessage(err, "failed to load timezone")
	}
	return loc, nil
}

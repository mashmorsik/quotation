package config

import (
	errs "github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Postgres struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"postgres"`
	Quotations []string `yaml:"quotations"`
	Server     struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	QuoteAPI struct {
		URL string `yaml:"url"`
	} `yaml:"quoteApi"`
	Cron struct {
		Location string `yaml:"location"`
		Period   string `yaml:"period"`
	} `yaml:"cron"`
	ResponseDelay time.Duration `yaml:"responseDelay"`
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errs.WithMessage(err, "failed to read config file")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errs.WithMessage(err, "failed to unmarshal config")
	}

	return &config, nil
}

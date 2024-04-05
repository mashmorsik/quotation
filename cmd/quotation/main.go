package main

import (
	"context"
	"github.com/mashmorsik/logger"
	"github.com/mashmorsik/quotation/config"
	cronSc "github.com/mashmorsik/quotation/cron"
	"github.com/mashmorsik/quotation/infrastructure/data"
	"github.com/mashmorsik/quotation/infrastructure/server"
	"github.com/mashmorsik/quotation/internal/quotation"
	"github.com/mashmorsik/quotation/repository"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.BuildLogger(nil)

	conf, err := config.LoadConfig()
	if err != nil {
		logger.Errf("Error loading config: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-sigCh
		logger.Infof("context done")
		cancel()
	}()

	conn := data.MustConnectPostgres(ctx, conf)
	data.MustMigrate(conn)
	dat := data.NewData(ctx, conn)

	quoteRepo := repository.NewQuoteRepo(ctx, dat)
	qq := quotation.NewQuotation(ctx, quoteRepo, conf)

	dt := cronSc.NewData(quoteRepo, conf)
	_, err = dt.RunScheduler()
	if err != nil {
		logger.Errf("Error running scheduler: %v", err)
		return
	}
	logger.Info("Scheduler started")

	httpServer := server.NewServer(conf, *qq)
	if err = httpServer.StartServer(); err != nil {
		logger.Err(err, "start http server failed")
	}
}

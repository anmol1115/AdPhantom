package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anmol1115/AdPhantom/internal/blocker"
	"github.com/anmol1115/AdPhantom/internal/config"
	Logger "github.com/anmol1115/AdPhantom/internal/logger"
	"github.com/anmol1115/AdPhantom/internal/resolver"
)

const (
	CONFIG_PATH  string = "/app/configs/config.ini"
	LOGFILE_PATH string = "/app/logs/app.log"
)

func main() {
	cfg, err := config.LoadConfig(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := Logger.Init(&cfg.Logging, LOGFILE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	err = blocker.LoadFilterRules(cfg.Filter.Lists)
	if err != nil {
		log.Fatal(err)
	}

	res := resolver.New(&cfg.DNS)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ctx = Logger.WithLogger(ctx, logger)
	ctx = resolver.WithResolver(ctx, res)

	go tcpListener(ctx)
	go udpListener(ctx)

	<-ctx.Done()
	logger.Debug("Exiting main thread")
}

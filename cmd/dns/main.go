package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anmol1115/AdPhantom/internal/config"
	Logger "github.com/anmol1115/AdPhantom/internal/logger"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ctx = Logger.WithLogger(ctx, logger)

	go tcpListener(ctx)
	go udpListener(ctx)

	<-ctx.Done()
	logger.Debug("Exiting main thread")
}

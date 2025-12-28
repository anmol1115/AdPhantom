package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anmol1115/AdPhantom/internal/config"
)

const CONFIG_PATH string = "/app/configs/config.ini"

func main() {
	_, err := config.LoadConfig(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go startConfigWatcher(ctx, CONFIG_PATH)
	go tcpListener()
	go udpListener()

	<-ctx.Done()
	log.Println("Exiting main thread")
}

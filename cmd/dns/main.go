package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anmol1115/AdPhantom/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("/app/configs/config.ini")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(cfg.Logging.Level) // remove: only for debugging

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go TcpListener()
	go UdpListener()

	<-ctx.Done()
	log.Println("Exiting main thread")
}

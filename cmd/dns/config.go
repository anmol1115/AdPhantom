package main

import (
	"context"
	"log"

	"github.com/anmol1115/AdPhantom/internal/config"
)

func startConfigWatcher(ctx context.Context, path string) {
	if err := config.Watch(ctx, path, onChange); err != nil {
		log.Println(err)
	}
}

func onChange(cfg *config.Config) {
	log.Printf("Updated config: %v", cfg)
}

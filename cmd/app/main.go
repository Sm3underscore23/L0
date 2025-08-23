package main

import (
	"flag"
	"log"

	"github.com/subosito/gotenv"
	"taskL0/internal/app"
	"taskL0/internal/config"
)

var (
	configPath = flag.String("config-path", "", "Path to config file")
	envLocal   = flag.Bool("env-local", false, "for .env load")
)

func main() {
	// Flags
	flag.Parse()

	// Env
	if *envLocal {
		if err := gotenv.Load(); err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	// Configuration
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}

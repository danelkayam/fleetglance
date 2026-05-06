package main

import (
	"flag"
	"fleetglance/internal/console"
	"fleetglance/internal/console/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Params struct {
	ConfigPath string `env:"FLEETGLANCE_CONFIG_PATH"`
	LogFormat  string `env:"LOG_FORMAT" envDefault:"console" validate:"oneof=console json"`
}

func main() {
	params, err := loadParams(os.Args[1:])
	if err != nil {
		fmt.Printf("Error loading parameters: %v\n", err)
		return
	}

	// init logger
	if params.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting fleetglance console...")

	fleet, err := config.LoadFleet(params.ConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed loading fleet config")
		return
	}

	c := console.NewConsole(fleet)

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(termChan)

	go func() {
		errChan <- c.Start()
	}()

	select {
	case sig := <-termChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Error().Err(err).Msg("Console stopped with error")
		}
	}

	log.Info().Msg("Shutting down fleetglance console...")

	if err := c.Stop(); err != nil {
		log.Error().Err(err).Msg("Failed stopping fleetglance console")
	}

	log.Info().Msg("Shutting down fleetglance console... DONE")
}

func loadParams(args []string) (*Params, error) {
	// Optional local development file.
	// Missing .env should not fail the service.
	_ = godotenv.Load(".env")

	params, err := env.ParseAs[Params]()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configPathFlag := flags.String("f", "", "path to fleet config file")
	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parse flags: %w", err)
	}

	flagProvided := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "f" {
			flagProvided = true
		}
	})

	if flagProvided {
		params.ConfigPath = *configPathFlag
	}

	if params.ConfigPath == "" {
		return nil, fmt.Errorf("fleet config path is required; set -f or FLEETGLANCE_CONFIG_PATH")
	}

	if err := validator.New().Struct(params); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	return &params, nil
}

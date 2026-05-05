package main

import (
	"fleetglance/internal/agent"
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
	Port      int    `env:"PORT" envDefault:"9800" validate:"min=1,max=65535"`
	Debug     bool   `env:"DEBUG" envDefault:"false"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"console" validate:"oneof=console json"`
}

func main() {
	params, err := loadParams()
	if err != nil {
		fmt.Printf("Error loading parameters: %v\n", err)
		return
	}

	// init logger
	if params.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting fleetglance agent...")

	agent := agent.NewAgent(agent.Params{
		Port:  params.Port,
		Debug: params.Debug,
	})

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errChan <- agent.Start()
	}()

	select {
	case sig := <-termChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Error().Err(err).Msg("Agent stopped with error")
		}
	}

	log.Info().Msg("Shutting down fleetglance agent...")

	if err := agent.Stop(); err != nil {
		log.Error().Err(err).Msg("Failed stopping fleetglance agent")
	}

	log.Info().Msg("Shutting down fleetglance agent... DONE")
}

func loadParams() (*Params, error) {
	// Optional local development file.
	// Missing .env should not fail the service.
	_ = godotenv.Load(".env")

	params, err := env.ParseAs[Params]()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	if err := validator.New().Struct(params); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	return &params, nil
}

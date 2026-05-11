package main

import (
	"context"
	"errors"
	"fleetglance/internal/console"
	"fleetglance/internal/console/config"
	"fleetglance/internal/version"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type Params struct {
	ConfigPath string `env:"FLEETGLANCE_CONFIG_PATH"`
	LogFormat  string `env:"LOG_FORMAT" envDefault:"console" validate:"oneof=console json"`
}

func main() {
	if err := newCommand(os.Stdout, os.Stderr, startConsole).Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newCommand(out io.Writer, errOut io.Writer, run func(*Params) error) *cli.Command {
	return &cli.Command{
		Name:        "fleetglance-console",
		Usage:       "run the Fleetglance console",
		Writer:      out,
		ErrWriter:   errOut,
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "version",
				Usage: "print version information and exit",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"f"},
				Usage:   "path to fleet config file",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Bool("version") {
				_, err := fmt.Fprint(out, version.Format("Console"))
				return err
			}

			params, err := resolveParams(cmd)
			if err != nil {
				return fmt.Errorf("load parameters: %w", err)
			}

			return run(params)
		},
	}
}

func startConsole(params *Params) error {
	// init logger
	if params.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting fleetglance console...")

	fleet, err := config.LoadFleet(params.ConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed loading fleet config")
		return fmt.Errorf("load fleet config: %w", err)
	}

	if err := config.ValidateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("Failed validating fleet config")
		return fmt.Errorf("validate fleet config: %w", err)
	}

	c := console.NewConsole(fleet)

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(termChan)

	go func() {
		errChan <- c.Start()
	}()

	var runErr error
	select {
	case sig := <-termChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Error().Err(err).Msg("Console stopped with error")
			runErr = fmt.Errorf("start console: %w", err)
		}
	}

	log.Info().Msg("Shutting down fleetglance console...")

	if err := c.Stop(); err != nil {
		log.Error().Err(err).Msg("Failed stopping fleetglance console")
		stopErr := fmt.Errorf("stop console: %w", err)
		if runErr != nil {
			return errors.Join(runErr, stopErr)
		}
		return stopErr
	}

	log.Info().Msg("Shutting down fleetglance console... DONE")

	return runErr
}

func resolveParams(cmd *cli.Command) (*Params, error) {
	params, err := loadParams()
	if err != nil {
		return nil, err
	}

	if cmd.IsSet("config") {
		params.ConfigPath = cmd.String("config")
	}

	if params.ConfigPath == "" {
		return nil, fmt.Errorf("fleet config path is required; set -f or FLEETGLANCE_CONFIG_PATH")
	}

	if err := validator.New().Struct(params); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	return params, nil
}

func loadParams() (*Params, error) {
	// Optional local development file.
	// Missing .env should not fail the service.
	_ = godotenv.Load(".env")

	params, err := env.ParseAs[Params]()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return &params, nil
}

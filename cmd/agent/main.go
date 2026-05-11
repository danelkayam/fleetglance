package main

import (
	"context"
	"errors"
	"fleetglance/internal/agent"
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
	Port      int    `env:"PORT" envDefault:"9800" validate:"min=1,max=65535"`
	Debug     bool   `env:"DEBUG" envDefault:"false"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"console" validate:"oneof=console json"`
}

func main() {
	if err := newCommand(os.Stdout, os.Stderr, startAgent).Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newCommand(out io.Writer, errOut io.Writer, run func(*Params) error) *cli.Command {
	return &cli.Command{
		Name:        "fleetglance-agent",
		Usage:       "run the Fleetglance telemetry agent",
		Writer:      out,
		ErrWriter:   errOut,
		HideVersion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "version",
				Usage: "print version information and exit",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Bool("version") {
				_, err := fmt.Fprint(out, version.Format("Agent"))
				return err
			}

			params, err := loadParams()
			if err != nil {
				return fmt.Errorf("load parameters: %w", err)
			}

			return run(params)
		},
	}
}

func startAgent(params *Params) error {
	// init logger
	if params.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting fleetglance agent...")

	a := agent.NewAgent(agent.Params{
		Port:  params.Port,
		Debug: params.Debug,
	})

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(termChan)

	go func() {
		errChan <- a.Start()
	}()

	var runErr error
	select {
	case sig := <-termChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-errChan:
		if err != nil {
			log.Error().Err(err).Msg("Agent stopped with error")
			runErr = fmt.Errorf("start agent: %w", err)
		}
	}

	log.Info().Msg("Shutting down fleetglance agent...")

	if err := a.Stop(); err != nil {
		log.Error().Err(err).Msg("Failed stopping fleetglance agent")
		stopErr := fmt.Errorf("stop agent: %w", err)
		if runErr != nil {
			return errors.Join(runErr, stopErr)
		}
		return stopErr
	}

	log.Info().Msg("Shutting down fleetglance agent... DONE")

	return runErr
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

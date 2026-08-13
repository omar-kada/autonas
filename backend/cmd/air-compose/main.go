// Package main is the entry point for AirCompose.
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"omar-kada/air-compose/internal/cli"
	"omar-kada/air-compose/internal/logs"
	"omar-kada/air-compose/internal/shell"

	"github.com/lmittmann/tint"
)

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	isDev := strings.ToUpper(os.Getenv("ENV")) == "DEV"
	var base slog.Handler
	if isDev {
		base =
			tint.NewTextHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.Kitchen,
				AddSource:  true,
			})

	} else {
		base = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logsHub := logs.NewHistoryHub(100)
	tapLogHandler := logs.NewTapHandler(base, func(_ context.Context, r slog.Record) {
		logsHub.Broadcast(r)
	})
	slog.SetDefault(slog.New(tapLogHandler))

	// Add subcommands
	rootCmd := cli.NewRootCmd(shell.NewExecutor(), logsHub)
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error executing root command", "error", err)
		retcode = 1 // it exits with code 1
	}
}

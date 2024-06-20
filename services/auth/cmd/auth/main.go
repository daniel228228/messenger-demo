package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	conf "messenger.auth/pkg/config"
	"messenger.auth/pkg/log"

	"messenger.auth/pkg/app"
	"messenger.auth/pkg/config"
)

const (
	GracefulShutdownTimeout = 60
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &config.Config{}

	err := conf.LoadConfig(config)
	if err != nil {
		panic(fmt.Sprintf("loading config error: %s", err.Error()))
	}

	logLevel, err := log.GetLevelFromString(config.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("setting logger error: %s", err.Error()))
	}

	stdout := log.NewHandler(os.Stdout, logLevel)

	log := log.NewLogger(app.Name).Add(stdout).Launch()

	log.Info().Bool("app", true).Str("state", "start").Send()
	defer log.Info().Bool("app", true).Str("state", "stop").Send()

	defer func() {
		if r := recover(); r != nil {
			log.Fatal().Bool("app", true).Msgf("panic: %v\n\n%s", r, string(debug.Stack()))
			cancel()
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		s := <-sig
		cause := s.String()

		log.Info().Bool("app", true).Str("state", "stopping").Str("cause", cause).Send()

		cancel()

		<-time.After(GracefulShutdownTimeout * time.Second)
		forceStop(log)
	}()

	app.Run(ctx, config, log)
}

func forceStop(log *log.Logger) {
	log.Warn().Bool("app", true).Msg("Failed to gracefully stop the current process, proceeding to force stop")
	log.Warn().Bool("app", true).Str("state", "stop").Bool("force", true).Send()

	os.Exit(1)
}

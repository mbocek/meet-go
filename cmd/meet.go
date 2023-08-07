package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/mbocek/meet-go/internal/config"
	"github.com/mbocek/meet-go/internal/logger"
	"github.com/mbocek/meet-go/internal/repository/postgres"
	"github.com/mbocek/meet-go/internal/route"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	const (
		console string = "console"
		json           = "json"
	)
	logFormat := flag.String("log", console, "Log output (console, json)")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// log message formatting
	switch *logFormat {
	case console:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	case json:
		logger.JsonLoggingSetup()
	}

	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	configuration := config.ReadConfigFile()
	r := gin.New()

	ctx := context.Background()
	dbRepo, errDB := postgres.New(ctx, configuration.Postgres)
	if errDB != nil {
		log.Fatal().Err(errDB).Msg("DB repo creation failed")
	}

	routes, errRoute := route.New(dbRepo)
	if errRoute != nil {
		log.Error().Err(errRoute).Msg("routes handler creation failed")
	}
	routes.RegisterHandlers(r)

	if err := r.Run(); err != nil {
		log.Fatal().Stack().Err(err).Msg("error running http server")
	}
}

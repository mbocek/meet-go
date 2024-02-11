package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/mbocek/meet-go/db"
	"github.com/mbocek/meet-go/internal/config"
	"github.com/mbocek/meet-go/internal/logger"
	"github.com/mbocek/meet-go/internal/repository/postgres"
	"github.com/mbocek/meet-go/internal/route/user"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var (
	demo                 *bool
	schemaMigrationsDemo = "schema_migrations_demo"
)

func init() {
	const (
		console string = "console"
		json           = "json"
	)
	logFormat := flag.String("log", console, "Log output (console, json)")
	debug := flag.Bool("debug", false, "Enable debug mode")
	demo = flag.Bool("demo", false, "Enable demo mode")
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

	if *demo {

	}
}

func migrateDB(url string, migrations embed.FS, migrationsPath string, migrationTable *string) error {
	if migrationTable != nil {
		log.Info().Str("table", *migrationTable).Msg("running db migration")
	} else {
		log.Info().Msg("running db migration")
	}
	source, err := iofs.New(migrations, migrationsPath)
	if err != nil {
		return eris.Wrap(err, "cannot open migrations")
	}

	var uri string
	if migrationTable != nil {
		// see https://github.com/golang-migrate/migrate/blob/master/database/mysql/README.md
		uri = fmt.Sprintf("%s&x-migrations-table=%s", url, *migrationTable)
	} else {
		uri = url
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, uri)
	if err != nil {
		return eris.Wrap(err, "cannot read migration")
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return eris.Wrap(err, "cannot migration all")
	}
	return nil
}

func main() {
	configuration := config.ReadConfigFile()
	r := gin.New()

	ctx := context.Background()
	dbRepo, errDB := postgres.New(ctx, configuration.Postgres)
	if errDB != nil {
		log.Fatal().Err(errDB).Msg("DB repo creation failed")
	}

	err := migrateDB(configuration.Postgres.Url, db.Migrations, "migrations", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot upgrade database")
	}
	if *demo {
		err := migrateDB(configuration.Postgres.Url, db.MigrationsDemo, "migrations-demo", &schemaMigrationsDemo)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot upgrade demo database")
		}
	}

	routes := user.NewRoute(dbRepo)
	routes.RegisterHandlers(r)

	if err := r.Run(); err != nil {
		log.Fatal().Stack().Err(err).Msg("error running http server")
	}
}

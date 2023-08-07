package test

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	"testing"
)

func Migrate(t *testing.T, url string, migrations embed.FS, migrationTable *string) {
	source, err := iofs.New(migrations, "migrations")
	require.NoError(t, err)

	var uri string
	if migrationTable != nil {
		// see https://github.com/golang-migrate/migrate/blob/master/database/mysql/README.md
		uri = fmt.Sprintf("%s&x-migrations-table=%s", url, *migrationTable)
	} else {
		uri = url
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, uri)
	require.NoError(t, err)
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}
}

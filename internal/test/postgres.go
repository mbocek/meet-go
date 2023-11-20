package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

type PostgresTestContainer struct {
	ctx        context.Context
	container  *postgres.PostgresContainer
	t          *testing.T
	credential credential
}

type credential struct {
	dbname   string
	user     string
	password string
}

func NewPostgres(t *testing.T, ctx context.Context) *PostgresTestContainer {
	pc := PostgresTestContainer{
		t:   t,
		ctx: ctx,
		credential: credential{
			dbname:   "meet",
			user:     "postgres",
			password: "password",
		},
	}
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16"),
		//postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		postgres.WithDatabase(pc.credential.dbname),
		postgres.WithUsername(pc.credential.user),
		postgres.WithPassword(pc.credential.password),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	pc.container = container
	return &pc
}

func (c *PostgresTestContainer) port() int {
	ctx, cancel := context.WithTimeout(c.ctx, time.Minute)
	defer cancel()
	p, err := c.container.MappedPort(ctx, "5432")
	require.NoError(c.t, err)
	return p.Int()
}

func (c *PostgresTestContainer) Url() string {
	return fmt.Sprintf("postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		c.credential.user,
		c.credential.password,
		c.port(),
		c.credential.dbname)
}

func (c *PostgresTestContainer) Close() {
	ctx, cancel := context.WithTimeout(c.ctx, time.Minute)
	defer cancel()
	err := c.container.Terminate(ctx)
	require.NoError(c.t, err)
}

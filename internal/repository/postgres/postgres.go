package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mbocek/meet-go/internal/config"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	db  *sqlx.DB
	ctx context.Context
}

func New(ctx context.Context, config config.Postgres) (*Repository, error) {
	log.Debug().Str("url", config.Url).Msg("opening DB connection")

	if db, errOpen := sqlx.Open("postgres", config.Url); errOpen != nil {
		return nil, errOpen
	} else {
		if errPing := db.Ping(); errPing != nil {
			return nil, errPing
		} else {
			log.Debug().Str("url", config.Url).Msg("opened DB connection")
			return &Repository{db, ctx}, nil
		}
	}
}

func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		log.Error().Err(err).Msg("error when closing postgres repo")
	}
}

func (r *Repository) GetDb() *sqlx.DB {
	return r.db
}

func (r *Repository) Ping() error {
	type DualData struct {
		Value int `db:"dual_count"`
	}

	query := `select 1 as dual_count from DUAL`

	var dualData DualData
	if err := r.GetDb().QueryRowx(query).StructScan(&dualData); err != nil {
		return err
	}

	if dualData.Value != 1 {
		return fmt.Errorf("ping returned wrong data %d", dualData.Value)
	}

	return nil
}

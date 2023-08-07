package interfaces

import "github.com/jmoiron/sqlx"

type DBRepository interface {
	Close()
	GetDb() *sqlx.DB
	Ping() error
}

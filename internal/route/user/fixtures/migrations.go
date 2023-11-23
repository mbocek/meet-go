package db

import "embed"

var (
	//go:embed migrations
	UserMigration embed.FS
)

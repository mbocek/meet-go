package db

import "embed"

//go:embed migrations
var Migrations embed.FS

//go:embed migrations-demo
var MigrationsDemo embed.FS

package config

import "github.com/elitekentoy/blog/internal/database"

type State struct {
	Config   *Config
	Database *database.Queries
}

package main

import (
	"context"
	"fmt"

	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
)

func middlewareLoggedIn(handler func(state *config.State, cmd config.Command, user *database.User) error) func(*config.State, config.Command) error {
	return func(state *config.State, command config.Command) error {
		user, err := state.Database.GetUser(context.Background(), state.Config.Username)
		if err != nil {
			return fmt.Errorf("error fetching user from the database: %w ", err)
		}

		return handler(state, command, &user)
	}

}

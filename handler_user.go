package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(state *config.State, command config.Command) error {
	if len(command.Arguments) != 1 {
		return fmt.Errorf("usage %s <name>", command.Name)
	}

	name := command.Arguments[0]

	user, err := state.Database.GetUser(context.Background(), name)
	if user.Name == "" || err != nil {
		os.Exit(1)
	}

	err = state.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched sucessfullly!")
	return nil
}

func handlerRegister(state *config.State, command config.Command) error {
	if len(command.Arguments) != 1 {
		return fmt.Errorf("usage %s <name>", command.Name)
	}

	name := command.Arguments[0]

	user, err := state.Database.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	state.Config.SetUser(name)
	fmt.Println(name, " has been created -> ", user)
	return nil
}

func handlerReset(state *config.State, command config.Command) error {
	err := state.Database.DeleteUser(context.Background())
	if err != nil {
		fmt.Println("database reset was failed")
		os.Exit(1)
	}

	fmt.Println("database reset was succesful")
	os.Exit(0)
	return nil
}

func handlerUsers(state *config.State, command config.Command) error {
	loggedInUser := state.Config.Username
	users, err := state.Database.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching users")
	}

	for _, user := range users {

		if user.Name == loggedInUser {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}

	}

	return nil
}

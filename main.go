package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf := config.ReadConfig()
	db, err := sql.Open("postgres", conf.DatabaseUrl)
	if err != nil {
		fmt.Printf("error occured when connecting to database")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	programState := config.State{
		Config:   &conf,
		Database: dbQueries,
	}
	commands := config.Commands{
		RegisteredCommands: make(map[string]func(*config.State, config.Command) error),
	}
	commands.Register("login", handlerLogin)
	commands.Register("register", handlerRegister)
	commands.Register("reset", handlerReset)
	commands.Register("users", handlerUsers)
	commands.Register("agg", handlerAgg)
	commands.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.Register("feeds", handlerFeeds)
	commands.Register("follow", middlewareLoggedIn(handlerFollow))
	commands.Register("following", middlewareLoggedIn(handlerFollowing))
	commands.Register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.Register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli<command> [args...]")
		return
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	err = commands.Run(&programState, &config.Command{Name: commandName, Arguments: commandArgs})
	if err != nil {
		log.Fatal(err)
	}

}

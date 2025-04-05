package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
)

func handlerBrowse(state *config.State, command config.Command, user *database.User) error {
	limit := 2

	if len(command.Arguments) == 1 {
		if specifiedLimit, err := strconv.Atoi(command.Arguments[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := state.Database.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}

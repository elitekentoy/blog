package main

import (
	"context"
	"fmt"
	"time"

	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(state *config.State, command config.Command, user *database.User) error {
	if len(command.Arguments) != 1 {
		return fmt.Errorf("usage: <url>")
	}

	url := command.Arguments[0]

	feed, err := state.Database.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed from the database: %w", err)
	}

	record, err := state.Database.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error inserting to feed follow database : %w", err)
	}

	fmt.Printf("\nFeed Name: %s", record.FeedName)
	fmt.Printf("\nUser: %s", record.UserName)

	return nil
}

func handlerFollowing(state *config.State, command config.Command, user *database.User) error {

	feeds, err := state.Database.GetFeedFolllowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching feeds from user: %w", err)
	}

	for _, feed := range feeds {
		println(feed.Name)
	}

	return nil
}

func handlerUnfollow(state *config.State, command config.Command, user *database.User) error {

	if len(command.Arguments) != 1 {
		return fmt.Errorf("usage: <url>")
	}

	url := command.Arguments[0]
	err := state.Database.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		Url:    url,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("error in unfollowing a feed: %w", err)
	}

	return nil
}

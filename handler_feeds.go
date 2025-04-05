package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elitekentoy/blog/internal/api/rss"
	"github.com/elitekentoy/blog/internal/config"
	"github.com/elitekentoy/blog/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(state *config.State, command config.Command) error {

	if len(command.Arguments) != 1 {
		return fmt.Errorf("usage: <time between requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(command.Arguments[0])
	if err != nil {
		return fmt.Errorf("invalid input for time between requests")
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(state)
	}
}

func scrapeFeeds(state *config.State) error {
	nextFeed, err := state.Database.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error in fetching feed: %w", err)
	}

	err = state.Database.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now()},
		ID:            nextFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error updating a feed to fetched: %w", err)
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error occured in fetching feed from api: %w", err)
	}

	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := state.Database.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    nextFeed.ID,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}

func handlerAddFeed(state *config.State, command config.Command, user *database.User) error {
	if len(command.Arguments) != 2 {
		os.Exit(1)
	}

	feedName := command.Arguments[0]
	feedUrl := command.Arguments[1]

	feed, err := state.Database.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error saving feed to database: %w", err)
	}
	fmt.Printf("feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	_, err = state.Database.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to followed feeds: %w", err)
	}
	return nil

}

func handlerFeeds(state *config.State, _ config.Command) error {

	feeds, err := state.Database.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching to database: %w", err)
	}

	userRecord := make(map[uuid.UUID]database.User)
	ids := []uuid.UUID{}

	for _, feed := range feeds {
		ids = append(ids, feed.UserID)
	}

	users, err := state.Database.GetUsersByID(context.Background(), ids)
	if err != nil {
		return fmt.Errorf("error fetching user from the database: %w", err)
	}

	for _, user := range users {
		userRecord[user.ID] = user
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", userRecord[feed.UserID].Name)

	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

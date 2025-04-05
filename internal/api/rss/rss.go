package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(context context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(context, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", res.Status)
	}

	decoder := xml.NewDecoder(res.Body)
	feed := &RSSFeed{}
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("error decoding feed: %w", err)
	}

	feed = unescapeStrings(feed)

	return feed, nil
}

func unescapeStrings(feed *RSSFeed) *RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return feed
}

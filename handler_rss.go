package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/inscrutabletaco/gator/internal/database"
)

const RSS_URL = "https://www.wagslane.dev/index.xml"

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println("Encountered an error scraping feeds:", err)
		}
	}

}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Feed created: %+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {

	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: feeds")
	}

	ctx := context.Background()

	results, err := s.db.GetFeedsByUser(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%-20s %-55s %-12s\n", "Feed Name", "Feed URL", "Owner")

	for i := 0; i < len(results); i++ {
		row := results[i]
		userName := row.Name_2.String
		if !row.Name_2.Valid {
			userName = "(unknown)"
		}
		fmt.Printf("%-20s %-55s %-12s\n", row.Name, row.Url, userName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	var params database.CreateFeedFollowParams
	params.UserID = user.ID
	params.FeedID = feed.ID

	result, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}

	fmt.Printf("Created feed follow: %v for %v", result.FeedName, result.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: following")
	}

	ctx := context.Background()

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, row := range feedFollows {
		fmt.Println(row.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	var params database.DeleteFeedFollowParams
	params.UserID = user.ID
	params.FeedID = feed.ID

	err = s.db.DeleteFeedFollow(ctx, params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed successfully unfollowed!")

	return nil
}

func scrapeFeeds(s *state) error {

	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to identify next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark as fetched feed %v: %w", nextFeed.Name, err)
	}

	rss, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed %v from %v: %w", nextFeed.Name, nextFeed.Url, err)
	}

	fmt.Printf("New items from feed %v\n", nextFeed.Name)

	for _, item := range rss.Channel.Item {
		fmt.Printf("   %v\n", item.Title)
	}

	return nil
}

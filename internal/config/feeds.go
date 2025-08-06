package config

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)
	if err != nil {
		fmt.Printf("error creating new HTTP request: %v\n", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error fetching RSS feed: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("status code %d returned from server", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading HTTP response: %v\n", err)
		return nil, err
	}

	feed := RSSFeed{}
	if err = xml.Unmarshal(data, &feed); err != nil {
		fmt.Printf("error unmarshalling bytes to RSSFeed struct: %v\n", err)
		return nil, err
	}

	unEscapeFields(&feed)
	return &feed, nil
}

func unEscapeFields(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
}

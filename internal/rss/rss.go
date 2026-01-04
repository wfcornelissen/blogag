package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		os.Exit(1)
		return &RSSFeed{}, fmt.Errorf("Failed to create client:\n%v\n", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		os.Exit(1)
		return &RSSFeed{}, fmt.Errorf("Request failed:\n%v\n", err)
	}
	res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Failed to read response body:\n%v\n", err)
	}

	var feed *RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Failed to unmarshal:\n%v\n", err)
	}

	cleanFeed := unescape(feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Failed to unescape:\n%v\n", err)
	}

	return cleanFeed, nil
}

func unescape(input *RSSFeed) *RSSFeed {
	var result *RSSFeed

	result.Channel.Title = html.UnescapeString(input.Channel.Title)
	result.Channel.Link = html.UnescapeString(input.Channel.Link)
	result.Channel.Description = html.UnescapeString(input.Channel.Description)

	var resultItem *RSSItem
	if len(result.Channel.Item) > 0 {
		for _, item := range input.Channel.Item {
			resultItem.Title = html.UnescapeString(item.Title)
			resultItem.Link = html.UnescapeString(item.Link)
			resultItem.Description = html.UnescapeString(item.Description)
			resultItem.PubDate = html.UnescapeString(item.PubDate)

			result.Channel.Item = append(result.Channel.Item, *resultItem)
		}
	}

	return result
}

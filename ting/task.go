package ting

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func RunTask(url string) error {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)

	if err != nil {
		return err
	}

	var urls []string
	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for _, item := range feed.Items {
		if strings.HasSuffix(item.Link, ".html") && startDate.Before(*item.PublishedParsed) {
			urls = append(urls, item.Link)
		}
	}

	return nil
}

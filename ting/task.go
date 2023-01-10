package ting

import (
	"github.com/mmcdole/gofeed"
	"log"
	"strings"
	"time"
)

func RunTask(url string) error {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)

	if err != nil {
		return err
	}

	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	var voas []*Voa

	for _, item := range feed.Items {
		if strings.HasSuffix(item.Link, ".html") && startDate.Before(*item.PublishedParsed) {
			voa, err := parseVoa(item.Link)

			if err != nil {
				log.Fatalf("Parse voa error %v", err)
			}

			voa.Title = item.Title
			voa.Description = item.Description
			voa.PublishedAtUtc = *item.PublishedParsed

			if len(item.Enclosures) > 0 {
				voa.ImageUrl = item.Enclosures[0].URL
			}

			voas = append(voas, voa)
		}
	}

	return nil
}

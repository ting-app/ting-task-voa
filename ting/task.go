package ting

import (
	"github.com/mmcdole/gofeed"
	"log"
	"strings"
	"time"
)

const programId = 2

func RunTask(channel Channel) error {
	url := channel.Url

	log.Printf("Start to fetch voa, url=%s\n", url)

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)

	if err != nil {
		return err
	}

	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var voas []*Voa

	for _, item := range feed.Items {
		if strings.HasSuffix(item.Link, ".html") && startDate.Before(*item.PublishedParsed) {
			voa, err := parseVoa(item.Link)

			if err != nil {
				log.Printf("failed to parse voa, url=%s", item.Link)

				return err
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

	if voas == nil || len(voas) == 0 {
		log.Printf("No news found, url=%s\n", url)

		return nil
	} else {
		log.Printf("Found %v news, url=%s\n", len(voas), url)
	}

	dbConfig, err := ParseDbConfig()

	if err != nil {
		return err
	}

	err = InitDb(dbConfig)

	if err != nil {
		return err
	}

	defer CloseDb()

	savedTing := 0

	for _, voa := range voas {
		ting := Ting{
			ProgramId:   programId,
			Title:       voa.Title,
			Description: voa.Description,
			AudioUrl:    voa.AudioUrl,
			Content:     voa.Body,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err = saveTing(ting)

		if err != nil {
			return err
		}

		savedTing += 1
	}

	log.Printf("Saved %v news as ting, url=%s\n", savedTing, url)

	return nil
}

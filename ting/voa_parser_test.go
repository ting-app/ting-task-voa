package ting

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestVoaParser(t *testing.T) {
	feedUrl := "https://learningenglish.voanews.com/api/zmg_pebmyp"
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedUrl)

	if err != nil {
		t.Log("parse feed error", err)
		t.Fail()
	}

	for _, item := range feed.Items {
		if strings.HasSuffix(item.Link, ".html") {
			voa, err := parseVoa(item.Link)

			if err != nil {
				t.Log("parse voa error", err)

				continue
			}

			assert.NotNil(t, voa)
			assert.NotNil(t, voa.Title)
			assert.NotNil(t, voa.Description)
			assert.NotNil(t, voa.Url)
			assert.NotNil(t, voa.PublishedAtUtc)
			assert.NotNil(t, voa.BodyWithHtml)
			assert.NotNil(t, voa.Body)
			assert.NotNil(t, voa.ImageUrl)
			assert.NotNil(t, voa.AudioUrl)
			assert.True(t, len(voa.Words) > 0)

			for _, word := range voa.Words {
				assert.NotNil(t, word.Word)
				assert.NotNil(t, word.PartOfSpeech)
				assert.NotNil(t, word.Definition)
			}
		}
	}
}

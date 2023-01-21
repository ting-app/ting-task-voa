package ting

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoaParser(t *testing.T) {
	feedUrl := "https://learningenglish.voanews.com/api/zpyp_e-rm_"
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedUrl)

	if err != nil {
		t.Log("parse feed error", err)
		t.Fail()
	}

	item := feed.Items[0]
	voa, err := parseVoa(item.Link)

	if err != nil {
		t.Log("parse voa error", err)
		t.Fail()
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
}

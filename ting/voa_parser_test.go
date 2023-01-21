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

	assert.NotNil(t, t, voa)
}

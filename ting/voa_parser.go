package ting

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type Voa struct {
	Title          string
	Description    string
	Url            string
	PublishedAtUtc time.Time
	Body           string
	ImageUrl       string
}

func parseVoa(url string) (*Voa, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code error: %d %s", response.StatusCode, response.Status))
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println(doc.Url)

	voa := &Voa{
		Url: url,
	}

	return voa, nil
}

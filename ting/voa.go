package ting

import "time"

type Voa struct {
	Title          string
	Description    string
	Url            string
	PublishedAtUtc time.Time
	Body           string
	BodyWithHtml   string
	ImageUrl       string
	AudioUrl       string
}

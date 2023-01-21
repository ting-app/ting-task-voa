package ting

import "time"

type Ting struct {
	ProgramId   int
	Title       string
	Description string
	AudioUrl    string
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

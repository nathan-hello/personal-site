package utils

import (
	"time"
)
type Image struct {
	Name     string
	Size string
        Ext string
	Url      string
	Alt  string
}

type Frontmatter struct {
	Author         string
	Title          string
	Date           time.Time
	Images         []Image
	Tags           []string
	OverrideHref   string
	OverrideLayout string
	Description    string
	Hidden         bool
}

type Blog struct {
	Id   int
        Url string
	Frnt Frontmatter
	Html string
}

type Comment struct {
        Id string
        Author string
        Date time.Time
        PostId string
        Text string
        Html string
}


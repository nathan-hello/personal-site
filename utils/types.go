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
        Comments []Comment
}

type Comment struct {
        Id int64
        Author string
        Date time.Time
        PostId int64
        Html string
}

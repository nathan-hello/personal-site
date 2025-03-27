package db

import (
	"time"

	"github.com/nathan-hello/personal-site/utils"
)

func (cmt Comment) NewBlogComment() utils.Comment {
	c := utils.Comment{}
	c.Id = cmt.ID
	c.Author = cmt.Author
	date, err := time.Parse(time.RFC3339, cmt.CreatedAt)
	if err != nil {
		date = time.Unix(0, 0)
	}
	c.Date = date
	c.PostId = cmt.PostID
	c.Html = cmt.Html
	return c
}

type UserRow = SelectUserByIdRow

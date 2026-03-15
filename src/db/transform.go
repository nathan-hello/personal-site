package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nathan-hello/personal-site/src/utils"
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
	if cmt.ImageID == nil {
		return c
	}
	image, err := Conn.SelectFromImage(context.Background(), *cmt.ImageID)
	if err != nil {
		return c
	}
	c.Image.Ext = image.Ext
	c.Image.Name = fmt.Sprintf("%s.%s", strconv.Itoa(int(image.ID)), c.Image.Ext)
	c.Image.BytesCount = image.Size
	c.Image.Size = utils.FormatSize(image.Size)
	c.Image.Url = fmt.Sprintf("/i/%s", c.Image.Name)

	return c
}

type MessagesByChatroom struct {
	Message Message
	Profile Profile
}

func GetMessagesByChatroom(ctx context.Context, arg SelectMessagesByChatroomParams) ([]MessagesByChatroom, error) {
	asdf, err := Conn.SelectMessagesByChatroom(ctx, arg)
	if err != nil {
		return nil, err
	}
	arr := []MessagesByChatroom{}
	for _, v := range asdf {
		arr = append(arr, MessagesByChatroom{
			Message: Message{
				ID:        v.ID,
				AuthorID:  v.AuthorID,
				Message:   v.Message,
				RoomID:    v.RoomID,
				CreatedAt: v.CreatedAt,
			},
			Profile: Profile{
				ID:       v.ID_2,
				Username: v.Username,
				Color:    v.Color,
			},
		})
	}
	return arr, nil
}

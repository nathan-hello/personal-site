package utils

import (
	"fmt"
	"time"
)

const TimeFormat = time.RFC3339


type ChatMessage struct {
	UserId    string    `json:"msg-userid"`
	Username  string    `json:"msg-author"`
	Text      string    `json:"msg-text"`
	Color     string    `json:"msg-color"`
	CreatedAt time.Time `json:"msg-time"`
}

// TimeToString(true)  = string(HH:MM)
//
// TimeToString(false) = .Format(time.RFC3339)
func (c *ChatMessage) TimeToString(hourMinute bool) string {
	if hourMinute {
		return fmt.Sprintf("%02d:%02d", c.CreatedAt.Hour(), c.CreatedAt.Minute())
	}
	return c.CreatedAt.Format(time.RFC3339)
}

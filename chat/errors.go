package chat

import "errors"

var (
        ErrNoTextInChatMsg = errors.New("illegal message - no text in chat message")
)


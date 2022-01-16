package functions

import (
	"context"
	"github.com/gotd/td/tg"
)

func GetMessages(context context.Context, client *tg.Client, messageIds []tg.InputMessageClass) (tg.MessageClassArray, error) {
	messages, err := client.MessagesGetMessages(context, messageIds)
	if err != nil {
		return nil, err
	}
	switch m := messages.(type) {
	case *tg.MessagesMessages:
		return m.Messages, nil
	case *tg.MessagesMessagesSlice:
		return m.Messages, nil
	case *tg.MessagesChannelMessages:
		return m.Messages, nil
	default:
		return nil, nil
	}
}

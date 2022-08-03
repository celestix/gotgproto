package functions

import (
	"context"

	"github.com/gotd/td/tg"
)

func GetChannelMessages(context context.Context, client *tg.Client, peer tg.InputChannelClass, messageIds []tg.InputMessageClass) (tg.MessageClassArray, error) {
	messages, err := client.ChannelsGetMessages(context, &tg.ChannelsGetMessagesRequest{
		Channel: peer,
		ID:      messageIds,
	})
	if err != nil {
		return nil, err
	}
	switch m := messages.(type) {
	case *tg.MessagesMessages:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	case *tg.MessagesMessagesSlice:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	case *tg.MessagesChannelMessages:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	default:
		return nil, nil
	}
}

func GetChatMessages(context context.Context, client *tg.Client, messageIds []tg.InputMessageClass) (tg.MessageClassArray, error) {
	messages, err := client.MessagesGetMessages(context, messageIds)
	if err != nil {
		return nil, err
	}
	switch m := messages.(type) {
	case *tg.MessagesMessages:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	case *tg.MessagesMessagesSlice:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	case *tg.MessagesChannelMessages:
		go SavePeersFromClassArray(m.Chats, m.Users)
		return m.Messages, nil
	default:
		return nil, nil
	}
}

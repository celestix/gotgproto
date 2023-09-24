package functions

import (
	"context"

	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

// GetChatIdFromPeer returns the chat/user id from the provided tg.PeerClass.
func GetChatIdFromPeer(peer tg.PeerClass) int64 {
	switch peer := peer.(type) {
	case *tg.PeerChannel:
		return peer.ChannelID
	case *tg.PeerUser:
		return peer.UserID
	case *tg.PeerChat:
		return peer.ChatID
	default:
		return 0
	}
}

// GetChatFromPeer returns the tg.ChatFull data of the provided tg.PeerClass.
func GetChatFromPeer(ctx context.Context, client *tg.Client, peer tg.PeerClass) (*tg.ChatFull, error) {
	switch peer := peer.(type) {
	case *tg.PeerChannel:
		chat, err := client.ChannelsGetFullChannel(ctx, &tg.InputChannel{
			ChannelID: peer.ChannelID,
		})
		if err != nil {
			return nil, err
		}
		return chat.FullChat.(*tg.ChatFull), nil
	case *tg.PeerChat:
		chat, err := client.MessagesGetFullChat(ctx, peer.ChatID)
		if err != nil {
			return nil, err
		}
		return chat.FullChat.(*tg.ChatFull), nil
	default:
		return nil, nil
	}
}

// GetInputPeerClassFromId finds provided user id in the session storage and returns it if found.
func GetInputPeerClassFromId(iD int64) tg.InputPeerClass {
	peer := storage.GetPeerById(iD)
	if peer.ID == 0 {
		return nil
	}
	switch storage.EntityType(peer.Type) {
	case storage.TypeUser:
		return &tg.InputPeerUser{
			UserID:     peer.ID,
			AccessHash: peer.AccessHash,
		}
	case storage.TypeChat:
		return &tg.InputPeerChat{
			ChatID: peer.ID,
		}
	case storage.TypeChannel:
		return &tg.InputPeerChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		}
	}
	return nil
}

func SavePeersFromClassArray(cs []tg.ChatClass, us []tg.UserClass) {
	for _, u := range us {
		u, ok := u.(*tg.User)
		if !ok {
			continue
		}
		storage.AddPeer(u.ID, u.AccessHash, storage.TypeUser, u.Username)
	}
	for _, c := range cs {
		switch c := c.(type) {
		case *tg.Channel:
			storage.AddPeer(c.ID, c.AccessHash, storage.TypeChannel, c.Username)
		case *tg.Chat:
			storage.AddPeer(c.ID, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
		}
	}
}

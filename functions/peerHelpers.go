package functions

import (
	"context"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func GetChatIdFromPeer(peer tg.PeerClass) int64 {
	switch peer.(type) {
	case *tg.PeerChannel:
		return peer.(*tg.PeerChannel).ChannelID
	case *tg.PeerUser:
		return peer.(*tg.PeerUser).UserID
	case *tg.PeerChat:
		return peer.(*tg.PeerChat).ChatID
	default:
		return 0
	}
}

func GetChatFromPeer(ctx context.Context, client *tg.Client, peer tg.PeerClass) (*tg.ChatFull, error) {
	switch peer.(type) {
	case *tg.PeerChannel:
		chat, err := client.ChannelsGetFullChannel(ctx, &tg.InputChannel{
			ChannelID: peer.(*tg.PeerChannel).ChannelID,
		})
		if err != nil {
			return nil, err
		}
		return chat.FullChat.(*tg.ChatFull), nil
	case *tg.PeerChat:
		chat, err := client.MessagesGetFullChat(ctx, peer.(*tg.PeerChat).ChatID)
		if err != nil {
			return nil, err
		}
		return chat.FullChat.(*tg.ChatFull), nil
	default:
		return nil, nil
	}
}

func GetInputPeerClassFromId(iD int64) tg.InputPeerClass {
	peer := storage.GetPeerById(iD)
	if peer.ID == 0 {
		return nil
	}
	switch peer.Type {
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

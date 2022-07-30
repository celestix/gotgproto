package functions

import (
	"context"
	"errors"

	"github.com/anonyindian/gotgproto/storage"
	"github.com/anonyindian/gotgproto/types"
	"github.com/gotd/td/tg"
)

func ExtractContactResolvedPeer(p *tg.ContactsResolvedPeer, err error) (types.EffectiveChat, error) {
	if err != nil {
		return &types.EmptyUC{}, err
	}
	go func() {
		for _, chat := range p.Chats {
			switch chat := chat.(type) {
			case *tg.Channel:
				storage.AddPeer(chat.ID, chat.AccessHash, storage.TypeChannel, chat.Username)
			}
		}
		for _, user := range p.Users {
			user, ok := user.(*tg.User)
			if !ok {
				continue
			}
			storage.AddPeer(user.ID, user.AccessHash, storage.TypeUser, user.Username)
		}
	}()
	switch p.Peer.(type) {
	case *tg.PeerChannel:
		if p.Chats == nil || len(p.Chats) == 0 {
			return &types.EmptyUC{}, errors.New("peer info not found in the resolved Chats")
		}
		switch chat := p.Chats[0].(type) {
		case *tg.Channel:
			var c = types.Channel(*chat)
			return &c, nil
		case *tg.ChannelForbidden:
			return &types.EmptyUC{}, errors.New("peer could not be resolved because Channel Forbidden")
		}
	case *tg.PeerUser:
		if p.Users == nil || len(p.Users) == 0 {
			return &types.EmptyUC{}, errors.New("peer info not found in the resolved Chats")
		}
		switch user := p.Users[0].(type) {
		case *tg.User:
			var c = types.User(*user)
			return &c, nil
		}
	}
	return &types.EmptyUC{}, errors.New("contact not found")
}

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

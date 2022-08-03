package functions

import (
	"context"

	"github.com/gotd/td/tg"
)

func AddChatMembers(context context.Context, client *tg.Client, chatPeer tg.InputPeerClass, users []tg.InputUserClass, forwardLimit int) (bool, error) {
	switch c := chatPeer.(type) {
	case *tg.InputPeerChat:
		for _, user := range users {
			user, ok := user.(*tg.InputUser)
			if ok {
				_, err := client.MessagesAddChatUser(context, &tg.MessagesAddChatUserRequest{
					ChatID: c.ChatID,
					UserID: &tg.InputUser{
						UserID:     user.UserID,
						AccessHash: user.AccessHash,
					},
					FwdLimit: forwardLimit,
				})
				if err != nil {
					return false, err
				}
			}
		}
		return true, nil
	case *tg.InputPeerChannel:
		_, err := client.ChannelsInviteToChannel(context, &tg.ChannelsInviteToChannelRequest{
			Channel: &tg.InputChannel{
				ChannelID:  c.ChannelID,
				AccessHash: c.AccessHash,
			},
			Users: users,
		})
		return err == nil, err
	}
	return false, nil
}

func ArchiveChats(context context.Context, client *tg.Client, peers []tg.InputPeerClass) (bool, error) {
	var folderPeers = make([]tg.InputFolderPeer, len(peers))
	for n, peer := range peers {
		folderPeers[n] = tg.InputFolderPeer{
			Peer:     peer,
			FolderID: 1,
		}
	}
	_, err := client.FoldersEditPeerFolders(context, folderPeers)
	return err == nil, err
}

func CreateChannel(context context.Context, client *tg.Client, title, about string, broadcast bool) (*tg.Channel, error) {
	udps, err := client.ChannelsCreateChannel(context, &tg.ChannelsCreateChannelRequest{
		Title:     title,
		About:     about,
		Broadcast: broadcast,
	})
	if err != nil {
		return nil, err
	}
	// Highly experimental value from ChatClass array
	_, chats, _ := getUpdateFromUpdates(udps)
	return chats[0].(*tg.Channel), nil
}

func CreateChat(context context.Context, client *tg.Client, title string, users []tg.InputUserClass) (*tg.Chat, error) {
	udps, err := client.MessagesCreateChat(context, &tg.MessagesCreateChatRequest{
		Users: users,
		Title: title,
	})
	if err != nil {
		return nil, err
	}
	// Highly experimental value from ChatClass map
	_, chats, _ := getUpdateFromUpdates(udps)
	return chats[0].(*tg.Chat), nil
}

func BanChatMember(context context.Context, client *tg.Client, chatPeer tg.InputPeerClass, userPeer *tg.InputPeerUser, untilDate int) (tg.UpdatesClass, error) {
	switch c := chatPeer.(type) {
	case *tg.InputPeerChannel:
		return client.ChannelsEditBanned(context, &tg.ChannelsEditBannedRequest{
			Channel: &tg.InputChannel{
				ChannelID:  c.ChannelID,
				AccessHash: c.AccessHash,
			},
			Participant: userPeer,
			BannedRights: tg.ChatBannedRights{
				UntilDate:    untilDate,
				ViewMessages: true,
				SendMessages: true,
				SendMedia:    true,
				SendStickers: true,
				SendGifs:     true,
				SendGames:    true,
				SendInline:   true,
				EmbedLinks:   true,
			},
		})
	case *tg.InputPeerChat:
		return client.MessagesDeleteChatUser(context, &tg.MessagesDeleteChatUserRequest{
			ChatID: c.ChatID,
			UserID: &tg.InputUser{
				UserID:     userPeer.UserID,
				AccessHash: userPeer.AccessHash,
			},
		})
	default:
		return &tg.Updates{}, nil
	}
}

func UnarchiveChats(context context.Context, client *tg.Client, peers []tg.InputPeerClass) (bool, error) {
	var folderPeers = make([]tg.InputFolderPeer, len(peers))
	for n, peer := range peers {
		folderPeers[n] = tg.InputFolderPeer{
			Peer:     peer,
			FolderID: 0,
		}
	}
	_, err := client.FoldersEditPeerFolders(context, folderPeers)
	return err == nil, err
}

func UnbanChatMember(context context.Context, client *tg.Client, chatPeer *tg.InputPeerChannel, userPeer *tg.InputPeerUser) (bool, error) {
	_, err := client.ChannelsEditBanned(context, &tg.ChannelsEditBannedRequest{
		Channel: &tg.InputChannel{
			ChannelID:  chatPeer.ChannelID,
			AccessHash: chatPeer.AccessHash,
		},
		Participant: userPeer,
		BannedRights: tg.ChatBannedRights{
			UntilDate: 0,
		},
	})
	return err == nil, err
}

package ext

import (
	"context"
	"strings"
	"time"

	"github.com/celestix/gotgproto/storage"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

// Update contains all the data related to an update.
type Update struct {
	// EffectiveMessage is the tg.Message of current update.
	EffectiveMessage *types.Message
	// CallbackQuery is the tg.UpdateBotCallbackQuery of current update.
	CallbackQuery *tg.UpdateBotCallbackQuery
	// InlineQuery is the tg.UpdateInlineBotCallbackQuery of current update.
	InlineQuery *tg.UpdateBotInlineQuery
	// ChatJoinRequest is the tg.UpdatePendingJoinRequests of current update.
	ChatJoinRequest *tg.UpdatePendingJoinRequests
	// ChatParticipant is the tg.UpdateChatParticipant of current update.
	ChatParticipant *tg.UpdateChatParticipant
	// ChannelParticipant is the tg.UpdateChannelParticipant of current update.
	ChannelParticipant *tg.UpdateChannelParticipant
	// UpdateClass is the current update in raw form.
	UpdateClass tg.UpdateClass
	// Entities of an update, i.e. mapped users, chats and channels.
	Entities *tg.Entities
}

// GetNewUpdate creates a new Update with provided parameters.
func GetNewUpdate(ctx context.Context, client *tg.Client, p *storage.PeerStorage, e *tg.Entities, update tg.UpdateClass) *Update {
	u := &Update{
		UpdateClass: update,
	}
	switch update := update.(type) {
	case *tg.UpdateNewMessage:
		m, ok := update.GetMessage().(*tg.Message)
		if ok {
			u.EffectiveMessage = types.ConstructMessage(m)
		}
		diff, err := client.UpdatesGetDifference(ctx, &tg.UpdatesGetDifferenceRequest{
			Pts:  update.Pts - 1,
			Date: int(time.Now().Unix()),
		})
		// Silently add catched entities to *tg.Entities
		if err == nil {
			if value, ok := diff.(*tg.UpdatesDifference); ok {
				for _, vu := range value.Chats {
					switch chat := vu.(type) {
					case *tg.Chat:
						go p.AddPeer(chat.ID, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
						e.Chats[chat.ID] = chat
					case *tg.Channel:
						go p.AddPeer(chat.ID, chat.AccessHash, storage.TypeChannel, chat.Username)
						e.Channels[chat.ID] = chat
					}
				}
				for _, vu := range value.Users {
					user, ok := vu.AsNotEmpty()
					if !ok {
						continue
					}
					go p.AddPeer(user.ID, user.AccessHash, storage.TypeUser, user.Username)
					e.Users[user.ID] = user
				}
			}
		}
	case message.AnswerableMessageUpdate:
		m, ok := update.GetMessage().(*tg.Message)
		if ok {
			u.EffectiveMessage = types.ConstructMessage(m)
		}
	case *tg.UpdateBotCallbackQuery:
		u.CallbackQuery = update
	case *tg.UpdateBotInlineQuery:
		u.InlineQuery = update
	case *tg.UpdatePendingJoinRequests:
		u.ChatJoinRequest = update
	case *tg.UpdateChatParticipant:
		u.ChatParticipant = update
	case *tg.UpdateChannelParticipant:
		u.ChannelParticipant = update
	}
	u.Entities = e
	return u
}

func (u *Update) Args() []string {
	switch {
	case u.EffectiveMessage != nil:
		return strings.Fields(u.EffectiveMessage.Text)
	case u.CallbackQuery != nil:
		return strings.Fields(string(u.CallbackQuery.Data))
	case u.InlineQuery != nil:
		return strings.Fields(u.InlineQuery.Query)
	default:
		return make([]string, 0)
	}
}

// EffectiveUser returns the tg.User who is responsible for the update.
func (u *Update) EffectiveUser() *tg.User {
	if u.Entities == nil {
		return nil
	}
	var userId int64
	switch {
	case u.EffectiveMessage != nil:
		uId, ok := u.EffectiveMessage.FromID.(*tg.PeerUser)
		if !ok {
			for _, user := range u.Entities.Users {
				if user.Self && user.Bot {
					return nil
				}
				return user
			}
		}
		userId = uId.UserID
	case u.CallbackQuery != nil:
		userId = u.CallbackQuery.UserID
	case u.InlineQuery != nil:
		userId = u.InlineQuery.UserID
	case u.ChatParticipant != nil:
		userId = u.ChannelParticipant.UserID
	case u.ChannelParticipant != nil:
		userId = u.ChannelParticipant.UserID
	}
	return u.Entities.Users[userId]
}

// GetChat returns the responsible tg.Chat for the current update.
func (u *Update) GetChat() *tg.Chat {
	if u.Entities == nil {
		return nil
	}
	var (
		peer tg.PeerClass
	)
	switch {
	case u.EffectiveMessage != nil:
		peer = u.EffectiveMessage.PeerID
	case u.CallbackQuery != nil:
		peer = u.CallbackQuery.Peer
	case u.ChatJoinRequest != nil:
		peer = u.ChatJoinRequest.Peer
	case u.ChatParticipant != nil:
		peer = &tg.PeerChat{ChatID: u.ChatParticipant.ChatID}
	}
	if peer == nil {
		return nil
	}
	c, ok := peer.(*tg.PeerChat)
	if !ok {
		return nil
	}
	return u.Entities.Chats[c.ChatID]
}

// GetChannel returns the responsible tg.Channel for the current update.
func (u *Update) GetChannel() *tg.Channel {
	if u.Entities == nil {
		return nil
	}
	var (
		peer tg.PeerClass
	)
	switch {
	case u.EffectiveMessage != nil:
		peer = u.EffectiveMessage.PeerID
	case u.CallbackQuery != nil:
		peer = u.CallbackQuery.Peer
	case u.ChatJoinRequest != nil:
		peer = u.ChatJoinRequest.Peer
	case u.ChannelParticipant != nil:
		peer = &tg.PeerChannel{ChannelID: u.ChannelParticipant.ChannelID}
	}
	if peer == nil {
		return nil
	}
	c, ok := peer.(*tg.PeerChannel)
	if !ok {
		return nil
	}
	return u.Entities.Channels[c.ChannelID]
}

// GetUserChat returns the responsible tg.User for the current update.
func (u *Update) GetUserChat() *tg.User {
	if u.Entities == nil {
		return nil
	}
	var (
		peer tg.PeerClass
	)
	switch {
	case u.EffectiveMessage != nil:
		peer = u.EffectiveMessage.PeerID
	case u.CallbackQuery != nil:
		peer = u.CallbackQuery.Peer
	case u.ChatJoinRequest != nil:
		peer = u.ChatJoinRequest.Peer
	case u.ChatParticipant != nil:
		peer = &tg.PeerChat{ChatID: u.ChatParticipant.ChatID}
	}
	if peer == nil {
		return nil
	}
	c, ok := peer.(*tg.PeerUser)
	if !ok {
		return nil
	}
	return u.Entities.Users[c.UserID]
}

// EffectiveChat returns the responsible EffectiveChat for the current update.
func (u *Update) EffectiveChat() types.EffectiveChat {
	if c := u.GetChannel(); c != nil {
		cn := types.Channel(*c)
		return &cn
	} else if c := u.GetChat(); c != nil {
		cn := types.Chat(*c)
		return &cn
	} else if c := u.GetUserChat(); c != nil {
		cn := types.User(*c)
		return &cn
	}
	return &types.EmptyUC{}
}

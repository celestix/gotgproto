package ext

import (
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

// Update contains all the data related to an update.
type Update struct {
	// EffectiveMessage is the tg.Message of current update.
	EffectiveMessage *tg.Message
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
func GetNewUpdate(e *tg.Entities, update tg.UpdateClass) *Update {
	u := &Update{
		UpdateClass: update,
		Entities:    e,
	}
	switch update := update.(type) {
	case message.AnswerableMessageUpdate:
		m, ok := update.GetMessage().(*tg.Message)
		if ok {
			u.EffectiveMessage = m
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
	return u
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

// EffectiveChat returns the responsible tg.ChatClass for the current update.
func (u *Update) EffectiveChat() tg.ChatClass {
	if c := u.GetChannel(); c != nil {
		return c
	} else if c := u.GetChat(); c != nil {
		return c
	}
	return &tg.ChatEmpty{}
}

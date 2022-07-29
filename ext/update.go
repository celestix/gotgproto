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

// EffectiveChat returns the responsible tg.ChatClass for the current update.
func (u *Update) EffectiveChat() tg.ChatClass {
	if c := u.GetChannel(); c != nil {
		return c
	} else if c := u.GetChat(); c != nil {
		return c
	}
	return &tg.ChatEmpty{}
}

func (u *Update) GetUnitedChat() UnitedChat {
	if c := u.GetChannel(); c != nil {
		cn := Channel(*c)
		return &cn
	} else if c := u.GetChat(); c != nil {
		cn := Chat(*c)
		return &cn
	} else if c := u.EffectiveUser(); c != nil {
		cn := User(*c)
		return &cn
	}
	return &EmptyUC{}
}

type UnitedChat interface {
	GetID() int64
	GetAccessHash() int64
	IsAChannel() bool
	IsAChat() bool
	IsAUser() bool
}

type EmptyUC struct{}

func (*EmptyUC) GetID() int64 {
	return 0
}
func (*EmptyUC) GetAccessHash() int64 {
	return 0
}
func (*EmptyUC) IsAChannel() bool {
	return false
}
func (*EmptyUC) IsAChat() bool {
	return false
}
func (*EmptyUC) IsAUser() bool {
	return false
}

type User tg.User

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) GetAccessHash() int64 {
	return u.AccessHash
}

func (u *User) IsAChannel() bool {
	return false
}

func (u *User) IsAChat() bool {
	return false
}

func (u *User) IsAUser() bool {
	return true
}

func (u *User) Raw() *tg.User {
	us := tg.User(*u)
	return &us
}

type Channel tg.Channel

func (u *Channel) GetID() int64 {
	return u.ID
}

func (u *Channel) GetAccessHash() int64 {
	return u.AccessHash
}

func (u *Channel) IsAChannel() bool {
	return true
}

func (u *Channel) IsAChat() bool {
	return false
}

func (u *Channel) IsAUser() bool {
	return false
}

func (u *Channel) Raw() *tg.Channel {
	us := tg.Channel(*u)
	return &us
}

type Chat tg.Chat

func (u *Chat) GetID() int64 {
	return u.ID
}

func (u *Chat) GetAccessHash() int64 {
	return 0
}

func (u *Chat) IsAChannel() bool {
	return true
}

func (u *Chat) IsAChat() bool {
	return false
}

func (u *Chat) IsAUser() bool {
	return false
}

func (u *Chat) Raw() *tg.Chat {
	us := tg.Chat(*u)
	return &us
}

package ext

import (
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

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
	case *tg.UpdateInlineBotCallbackQuery:
		u.InlineQuery = update
	}
	return u
}

type Update struct {
	EffectiveMessage *tg.Message
	CallbackQuery    *tg.UpdateBotCallbackQuery
	InlineQuery      *tg.UpdateInlineBotCallbackQuery
	UpdateClass      tg.UpdateClass
	Entities         *tg.Entities
}

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

func (u *Update) EffectiveChat() tg.ChatClass {
	if c := u.GetChannel(); c != nil {
		return c
	} else if c := u.GetChat(); c != nil {
		return c
	}
	return &tg.ChatEmpty{}
}

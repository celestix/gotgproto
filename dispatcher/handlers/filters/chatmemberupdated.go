package filters

import "github.com/KoNekoD/gotgproto/ext"

type chatMemberUpdated struct{}

// All returns true on every type of tg.UpdateChannelParticipant and tg.UpdateChatParticipant update.
func (*chatMemberUpdated) All(_ *ext.Update) bool {
	return true
}

// ChatUpdate returns true on every type of tg.UpdateChatParticipant update.
func (*chatMemberUpdated) ChatUpdate(u *ext.Update) bool {
	return u.ChatParticipant != nil
}

// ChannelUpdate returns true on every type of tg.UpdateChannelParticipant update.
func (*chatMemberUpdated) ChannelUpdate(u *ext.Update) bool {
	return u.ChannelParticipant != nil
}

// FromUserId checks if the tg.UpdateChatParticipant and tg.UpdateChannelParticipant was sent by the provided user id and returns true if matches.
func (*chatMemberUpdated) FromUserId(userId int64) ChatMemberUpdatedFilter {
	return func(u *ext.Update) bool {
		if u.ChannelParticipant != nil {
			return u.ChannelParticipant.UserID == userId
		}
		if u.ChatParticipant != nil {
			return u.ChatParticipant.UserID == userId
		}
		return false
	}
}

// FromChatId checks if the tg.UpdateChatParticipant and tg.UpdateChannelParticipant was sent at the provided chat id and returns true if matches.
func (*chatMemberUpdated) FromChatId(chatId int64) ChatMemberUpdatedFilter {
	return func(u *ext.Update) bool {
		if u.ChannelParticipant != nil {
			return u.ChannelParticipant.ChannelID == chatId
		}
		if u.ChatParticipant != nil {
			return u.ChatParticipant.ChatID == chatId
		}
		return false
	}
}

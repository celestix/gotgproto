package filters

import (
	"github.com/gotd/td/tg"
	"strings"
)

type inlineQuery struct{}

// All returns true on every type of tg.UpdateBotInlineQuery update.
func (*inlineQuery) All(iq *tg.UpdateBotInlineQuery) bool {
	return true
}

// Prefix returns true if the tg.UpdateBotInlineQuery's Query field contains provided prefix.
func (*inlineQuery) Prefix(prefix string) InlineQueryFilter {
	return func(iq *tg.UpdateBotInlineQuery) bool {
		return strings.HasPrefix(iq.Query, prefix)
	}
}

// Suffix returns true if the tg.UpdateBotInlineQuery's Query field contains provided suffix.
func (*inlineQuery) Suffix(suffix string) InlineQueryFilter {
	return func(iq *tg.UpdateBotInlineQuery) bool {
		return strings.HasPrefix(iq.Query, suffix)
	}
}

// Equal checks if the tg.UpdateBotInlineQuery's Query field is equal to the provided data and returns true if matches.
func (*inlineQuery) Equal(data string) InlineQueryFilter {
	return func(iq *tg.UpdateBotInlineQuery) bool {
		return iq.Query == data
	}
}

// FromUserId checks if the tg.UpdateBotInlineQuery was sent by the provided user id and returns true if matches.
func (*inlineQuery) FromUserId(userId int64) InlineQueryFilter {
	return func(iq *tg.UpdateBotInlineQuery) bool {
		return iq.UserID == userId
	}
}

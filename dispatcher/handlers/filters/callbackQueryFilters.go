package filters

import (
	"github.com/gotd/td/tg"
	"strings"
)

// CallbackQueryPrefix checks if the tg.UpdateBotCallbackQuery's Data field has provided prefix.
func CallbackQueryPrefix(prefix string) CallbackQueryFilter {
	return func(cbq *tg.UpdateBotCallbackQuery) bool {
		return strings.HasPrefix(string(cbq.Data), prefix)
	}
}

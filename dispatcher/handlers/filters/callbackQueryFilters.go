package filters

import (
	"github.com/gotd/td/tg"
	"strings"
)

func CallbackQueryPrefix(prefix string) CallbackQueryFilter {
	return func(cbq *tg.UpdateBotCallbackQuery) bool {
		return strings.HasPrefix(string(cbq.Data), prefix)
	}
}

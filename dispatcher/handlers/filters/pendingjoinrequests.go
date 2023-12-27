package filters

import (
	"github.com/celestix/gotgproto/functions"
	"github.com/gotd/td/tg"
)

type pendingJoinRequests struct{}

// All returns true on every type of tg.UpdatePendingJoinRequests update.
func (*pendingJoinRequests) All(_ *tg.UpdatePendingJoinRequests) bool {
	return true
}

// ChatID returns true if the unique identifier of the chat where the UpdatePendingJoinRequests was created matches the input chatId.
func (*pendingJoinRequests) ChatID(chatId int64) PendingJoinRequestsFilter {
	return func(cjr *tg.UpdatePendingJoinRequests) bool {
		return functions.GetChatIdFromPeer(cjr.Peer) == chatId
	}
}

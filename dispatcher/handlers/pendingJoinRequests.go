package handlers

import (
	"github.com/amupxm/gotgproto/dispatcher/handlers/filters"
	"github.com/amupxm/gotgproto/ext"
)

// PendingJoinRequests handler is executed on all type of incoming updates.
type PendingJoinRequests struct {
	Callback CallbackResponse
	Filters  filters.PendingJoinRequestsFilter
}

// NewChatJoinRequest creates a new AnyUpdate handler bound to call its response.
func NewChatJoinRequest(filters filters.PendingJoinRequestsFilter, response CallbackResponse) PendingJoinRequests {
	return PendingJoinRequests{Callback: response, Filters: filters}
}

func (c PendingJoinRequests) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.ChatJoinRequest == nil {
		return nil
	}
	if c.Filters != nil && !c.Filters(u.ChatJoinRequest) {
		return nil
	}
	return c.Callback(ctx, u)
}

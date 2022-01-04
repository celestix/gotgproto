package dispatcher

import (
	"github.com/anonyindian/gotgproto/ext"
	"sort"
)

// Handler is the common interface for all the handlers.
type Handler interface {
	// CheckUpdate checks whether the update should be handled by this handler and processes it.
	CheckUpdate(*ext.Context, *ext.Update) error
}

// AddHandler adds a new handler to the dispatcher. The dispatcher will call CheckUpdate() to see whether the handler
// should be executed, and then execute it.
func (dp *CustomDispatcher) AddHandler(h Handler) {
	dp.AddHandlerToGroup(h, 0)
}

// AddHandlerToGroup adds a handler to a specific group; lowest number will be processed first.
func (dp *CustomDispatcher) AddHandlerToGroup(h Handler, group int) {
	handlers, ok := dp.handlerMap[group]
	if !ok {
		dp.handlerGroups = append(dp.handlerGroups, group)
		sort.Ints(dp.handlerGroups)
	}
	dp.handlerMap[group] = append(handlers, h)
}

package dispatcher

import (
	"github.com/anonyindian/gotgproto/ext"
	"sort"
)

type Handler interface {
	CheckUpdate(*ext.Context, *ext.Update) error
}

func (dp *CustomDispatcher) AddHandler(h Handler) {
	dp.AddHandlerToGroup(h, 0)
}

func (dp *CustomDispatcher) AddHandlerToGroup(h Handler, group int) {
	handlers, ok := dp.handlerMap[group]
	if !ok {
		dp.handlerGroups = append(dp.handlerGroups, group)
		sort.Ints(dp.handlerGroups)
	}
	dp.handlerMap[group] = append(handlers, h)
}

package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
	"log"
	"runtime/debug"
)

var (
	EndGroups        = errors.New("stopped")
	ContinueGroups   = errors.New("continued")
	SkipCurrentGroup = errors.New("skipped")
)

type CustomDispatcher struct {
	Panic         PanicHandler
	Error         ErrorHandler
	handlerMap    map[int][]Handler
	handlerGroups []int
}

type PanicHandler func(*ext.Context, *ext.Update, string)
type ErrorHandler func(*ext.Context, *ext.Update, string) error

func MakeDispatcher() *CustomDispatcher {
	return &CustomDispatcher{
		handlerMap: make(map[int][]Handler),
	}
}

type entities tg.Entities

func (u *entities) short() {
	u.Short = true
	u.Users = make(map[int64]*tg.User, 0)
	u.Chats = make(map[int64]*tg.Chat, 0)
	u.Channels = make(map[int64]*tg.Channel, 0)
}

func (dp *CustomDispatcher) Handle(ctx context.Context, updates tg.UpdatesClass) error {
	var (
		e    entities
		upds []tg.UpdateClass
	)
	switch u := updates.(type) {
	case *tg.Updates:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
		go func() {
			saveUsersPeers(u.Users)
			saveChatsPeers(u.Chats)
		}()
	case *tg.UpdatesCombined:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
		go func() {
			saveUsersPeers(u.Users)
			saveChatsPeers(u.Chats)
		}()
	case *tg.UpdateShort:
		upds = []tg.UpdateClass{u.Update}
		e.short()
	default:
		// *UpdateShortMessage
		// *UpdateShortChatMessage
		// *UpdateShortSentMessage
		// *UpdatesTooLong
		return nil
	}

	var err error
	for _, update := range upds {
		multierr.AppendInto(&err, dp.dispatch(ctx, tg.Entities(e), update))
	}
	return err
}

func (dp *CustomDispatcher) dispatch(ctx context.Context, e tg.Entities, update tg.UpdateClass) error {
	if update == nil {
		return nil
	}
	return dp.handleUpdates(ctx, e, update)
}

func (dp *CustomDispatcher) handleUpdates(ctx context.Context, e tg.Entities, update tg.UpdateClass) error {
	c := ext.NewContext(ctx, gotgproto.Api, gotgproto.Self, gotgproto.Sender, &e)
	u := ext.GetNewUpdate(&e, update)
	defer func() {
		if r := recover(); r != nil {
			errorStack := fmt.Sprintf("%s\n", r) + string(debug.Stack())
			if dp.Panic != nil {
				dp.Panic(c, u, errorStack)
				return
			}
			log.Println(errorStack)
		}
	}()
	var err error
	for group := range dp.handlerGroups {
		for _, handler := range dp.handlerMap[group] {
			err = handler.CheckUpdate(c, u)
			if err == nil || errors.Is(err, ContinueGroups) {
				continue
			} else if errors.Is(err, EndGroups) {
				return err
			} else if errors.Is(err, SkipCurrentGroup) {
				break
			} else {
				err = dp.Error(c, u, err.Error())
				switch err {
				case ContinueGroups:
					continue
				case EndGroups:
					return err
				case SkipCurrentGroup:
					break
				}
			}
		}
	}
	return err
}

func saveUsersPeers(u tg.UserClassArray) {
	for _, user := range u {
		c, ok := user.AsNotEmpty()
		if !ok {
			continue
		}
		storage.AddPeer(c.ID, c.AccessHash, storage.TypeUser, c.Username)
	}
}

func saveChatsPeers(u tg.ChatClassArray) {
	for _, chat := range u {
		channel, ok := chat.(*tg.Channel)
		if ok {
			storage.AddPeer(channel.ID, channel.AccessHash, storage.TypeChannel, channel.Username)
			continue
		}
		chat, ok := chat.(*tg.Chat)
		if !ok {
			continue
		}
		storage.AddPeer(chat.ID, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
	}
}

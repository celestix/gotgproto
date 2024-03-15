package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
)

var (
	// StopClient cancels the context and stops the client if returned through handler callback function.
	StopClient = errors.New("disconnect")

	// EndGroups stops iterating over handlers groups if returned through handler callback function.
	EndGroups = errors.New("stopped")
	// ContinueGroups continues iterating over handlers groups if returned through handler callback function.
	ContinueGroups = errors.New("continued")
	// SkipCurrentGroup skips current group and continues iterating over handlers groups if returned through handler callback function.
	SkipCurrentGroup = errors.New("skipped")
)

type Dispatcher interface {
	Initialize(context.Context, context.CancelFunc, *telegram.Client, *tg.User)
	Handle(context.Context, tg.UpdatesClass) error
	AddHandler(Handler)
	AddHandlerToGroup(Handler, int)
}

type NativeDispatcher struct {
	cancel              context.CancelFunc
	client              *tg.Client
	self                *tg.User
	sender              *message.Sender
	setReply            bool
	setEntireReplyChain bool
	// Panic handles all the panics that occur during handler execution.
	Panic PanicHandler
	// Error handles all the unknown errors which are returned by the handler callback functions.
	Error ErrorHandler
	// handlerMap is used for internal functionality of NativeDispatcher.
	handlerMap map[int][]Handler
	// handlerGroups is used for internal functionality of NativeDispatcher.
	handlerGroups []int

	pStorage *storage.PeerStorage
}

type PanicHandler func(*ext.Context, *ext.Update, string)
type ErrorHandler func(*ext.Context, *ext.Update, string) error

// MakeDispatcher creates new custom dispatcher which process and handles incoming updates.
func NewNativeDispatcher(setReply bool, setEntireReplyChain bool, eHandler ErrorHandler, pHandler PanicHandler, p *storage.PeerStorage) *NativeDispatcher {
	if eHandler == nil {
		eHandler = defaultErrorHandler
	}
	return &NativeDispatcher{
		pStorage:            p,
		handlerMap:          make(map[int][]Handler),
		handlerGroups:       make([]int, 0),
		setReply:            setReply,
		setEntireReplyChain: setEntireReplyChain,
		Error:               eHandler,
		Panic:               pHandler,
	}
}

func defaultErrorHandler(_ *ext.Context, _ *ext.Update, err string) error {
	log.Println("An error occured while handling update:", err)
	return ContinueGroups
}

type entities tg.Entities

func (u *entities) short() {
	u.Short = true
	u.Users = make(map[int64]*tg.User, 0)
	u.Chats = make(map[int64]*tg.Chat, 0)
	u.Channels = make(map[int64]*tg.Channel, 0)
}

func (dp *NativeDispatcher) Initialize(ctx context.Context, cancel context.CancelFunc, client *telegram.Client, self *tg.User) {
	dp.client = client.API()
	dp.sender = message.NewSender(dp.client)
	dp.self = self
	dp.cancel = cancel
}

// Handle function handles all the incoming updates, map entities and dispatches updates for further handling.
func (dp *NativeDispatcher) Handle(ctx context.Context, updates tg.UpdatesClass) error {
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
		go saveUsersPeers(u.Users, dp.pStorage)
		go saveChatsPeers(u.Chats, dp.pStorage)
	case *tg.UpdatesCombined:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
		go saveUsersPeers(u.Users, dp.pStorage)
		go saveChatsPeers(u.Chats, dp.pStorage)
	case *tg.UpdateShort:
		upds = []tg.UpdateClass{u.Update}
		e.short()
	default:
		return nil
	}

	var err error
	for _, update := range upds {
		multierr.AppendInto(&err, dp.dispatch(ctx, tg.Entities(e), update))
	}
	return err
}

func (dp *NativeDispatcher) dispatch(ctx context.Context, e tg.Entities, update tg.UpdateClass) error {
	if update == nil {
		return nil
	}
	return dp.handleUpdate(ctx, e, update)
}

func (dp *NativeDispatcher) handleUpdate(ctx context.Context, e tg.Entities, update tg.UpdateClass) error {
	u := ext.GetNewUpdate(ctx, dp.client, dp.pStorage, &e, update)
	dp.handleUpdateRepliedToMessage(u, ctx)
	c := ext.NewContext(ctx, dp.client, dp.pStorage, dp.self, dp.sender, &e, dp.setReply)
	var err error
	defer func() {
		if r := recover(); r != nil {
			errorStack := fmt.Sprintf("%s\n", r) + string(debug.Stack())
			if dp.Panic != nil {
				dp.Panic(c, u, errorStack)
				return
			} else {
				log.Println(errorStack)
			}
		}
	}()
	for _, group := range dp.handlerGroups {
		for _, handler := range dp.handlerMap[group] {
			err = handler.CheckUpdate(c, u)
			if err == nil || errors.Is(err, ContinueGroups) {
				continue
			} else if errors.Is(err, EndGroups) {
				return err
			} else if errors.Is(err, SkipCurrentGroup) {
				break
			} else if errors.Is(err, StopClient) {
				dp.cancel()
				return nil
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

func (dp *NativeDispatcher) handleUpdateRepliedToMessage(u *ext.Update, ctx context.Context) {
	msg := u.EffectiveMessage
	if msg == nil || !dp.setReply {
		return
	}
	for {
		if msg.Message.ReplyTo == nil {
			return
		}

		_ = msg.SetRepliedToMessage(ctx, dp.client, dp.pStorage)
		if !dp.setEntireReplyChain {
			return
		}
		msg = msg.ReplyToMessage
	}
}

func saveUsersPeers(u tg.UserClassArray, p *storage.PeerStorage) {
	for _, user := range u {
		c, ok := user.AsNotEmpty()
		if !ok {
			continue
		}
		p.AddPeer(c.ID, c.AccessHash, storage.TypeUser, c.Username)
	}
}

func saveChatsPeers(u tg.ChatClassArray, p *storage.PeerStorage) {
	for _, chat := range u {
		channel, ok := chat.(*tg.Channel)
		if ok {
			p.AddPeer(channel.ID, channel.AccessHash, storage.TypeChannel, channel.Username)
			continue
		}
		chat, ok := chat.(*tg.Chat)
		if !ok {
			continue
		}
		p.AddPeer(chat.ID, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
	}
}

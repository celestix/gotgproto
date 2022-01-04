package ext

import (
	"context"
	"github.com/anonyindian/gotgproto/functions"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
	"time"
)

type Context struct {
	OriginContext context.Context
	Client        *tg.Client
	Self          *tg.User
	Sender        *message.Sender
	Entities      *tg.Entities
}

func NewContext(ctx context.Context, client *tg.Client, self *tg.User, sender *message.Sender, entities *tg.Entities) *Context {
	return &Context{
		OriginContext: ctx,
		Client:        client,
		Self:          self,
		Sender:        sender,
		Entities:      entities,
	}
}

type ReplyOpts struct {
	NoWebpage        bool
	Markup           tg.ReplyMarkupClass
	ReplyToMessageId int
}

func (ctx *Context) Reply(upd *Update, text interface{}, opts *ReplyOpts) (tg.UpdatesClass, error) {
	if text == nil {
		return nil, ErrTextEmpty
	}
	if opts == nil {
		opts = &ReplyOpts{}
	}
	builder := ctx.Sender.Reply(*ctx.Entities, upd.UpdateClass.(message.AnswerableMessageUpdate))
	if opts.NoWebpage {
		builder = builder.NoWebpage()
	}
	if opts.Markup != nil {
		builder = builder.Markup(opts.Markup)
	}
	if opts.ReplyToMessageId != 0 {
		builder = builder.Reply(opts.ReplyToMessageId)
	}
	switch text := (text).(type) {
	case string:
		return builder.Text(ctx.OriginContext, text)
	case []styling.StyledTextOption:
		return builder.StyledText(ctx.OriginContext, text...)
	default:
		return nil, ErrTextInvalid
	}
}

func (ctx *Context) SendMessage(chatId int64, request tg.MessagesSendMessageRequest) (tg.UpdatesClass, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendMessage(ctx.OriginContext, &request)
}

func (ctx *Context) SendMedia(chatId int64, request *tg.MessagesSendMediaRequest) (tg.UpdatesClass, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendMedia(ctx.OriginContext, request)
}

func (ctx *Context) SendInlineBotResult(chatId int64, request *tg.MessagesSendInlineBotResultRequest) (tg.UpdatesClass, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendInlineBotResult(ctx.OriginContext, request)
}

func (ctx *Context) SendReaction(chatId int64, request *tg.MessagesSendReactionRequest) (tg.UpdatesClass, error) {
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendReaction(ctx.OriginContext, request)
}

func (ctx *Context) SendMultiMedia(chatId int64, request *tg.MessagesSendMultiMediaRequest) (tg.UpdatesClass, error) {
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendMultiMedia(ctx.OriginContext, request)
}

func (ctx *Context) GetChat(chatId int64) (tg.ChatFullClass, error) {
	peer := storage.GetPeerById(chatId)
	if peer.ID == 0 {
		return nil, ErrPeerNotFound
	}
	switch peer.Type {
	case storage.TypeChannel:
		channel, err := ctx.Client.ChannelsGetFullChannel(ctx.OriginContext, &tg.InputChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return channel.FullChat, nil
	case storage.TypeChat:
		chat, err := ctx.Client.MessagesGetFullChat(ctx.OriginContext, chatId)
		if err != nil {
			return nil, err
		}
		return chat.FullChat, nil
	}
	return nil, ErrNotChat
}

func (ctx *Context) GetUser(userId int64) (*tg.UserFull, error) {
	peer := storage.GetPeerById(userId)
	if peer.ID == 0 {
		return nil, ErrPeerNotFound
	}
	if peer.Type == storage.TypeUser {
		user, err := ctx.Client.UsersGetFullUser(ctx.OriginContext, &tg.InputUser{
			UserID:     peer.ID,
			AccessHash: peer.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return &user.FullUser, nil
	} else {
		return nil, ErrNotUser
	}
}

func (ctx *Context) ExportSessionString() (string, error) {
	return functions.EncodeSessionToString(storage.GetSession())
}

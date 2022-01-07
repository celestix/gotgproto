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

// Context consists of context.Context, tg.Client, Self etc.
type Context struct {
	// original context of an update.
	OriginContext context.Context
	// tg client which will be used make send requests.
	Client *tg.Client
	// self user who authorized the session.
	Self *tg.User
	// Sender is a message sending helper.
	Sender *message.Sender
	// Entities consists of mapped users, chats and channels from the update.
	Entities *tg.Entities
}

// NewContext creates a new Context object with provided parameters.
func NewContext(ctx context.Context, client *tg.Client, self *tg.User, sender *message.Sender, entities *tg.Entities) *Context {
	return &Context{
		OriginContext: ctx,
		Client:        client,
		Self:          self,
		Sender:        sender,
		Entities:      entities,
	}
}

// ReplyOpts object contains optional parameters for Context.Reply.
type ReplyOpts struct {
	// Whether the message should show link preview or not.
	NoWebpage bool
	// Reply markup of a message, i.e. inline keyboard buttons etc.
	Markup           tg.ReplyMarkupClass
	ReplyToMessageId int
}

// Reply uses given message update to create message for same chat and create a reply.
// Parameter 'text' interface should be one from string or an array of styling.StyledTextOption.
func (ctx *Context) Reply(upd *Update, text interface{}, opts *ReplyOpts) (*tg.Message, error) {
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
		return functions.ReturnNewMessageWithError(builder.Text(ctx.OriginContext, text))
	case []styling.StyledTextOption:
		return functions.ReturnNewMessageWithError(builder.StyledText(ctx.OriginContext, text...))
	default:
		return nil, ErrTextInvalid
	}
}

// SendMessage invokes method messages.sendMessage#d9d75a4 returning error if any.
func (ctx *Context) SendMessage(chatId int64, request tg.MessagesSendMessageRequest) (*tg.Message, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnNewMessageWithError(ctx.Client.MessagesSendMessage(ctx.OriginContext, &request))
}

// SendMedia invokes method messages.sendMedia#e25ff8e0 returning error if any. Send a media
func (ctx *Context) SendMedia(chatId int64, request tg.MessagesSendMediaRequest) (*tg.Message, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnNewMessageWithError(ctx.Client.MessagesSendMedia(ctx.OriginContext, &request))
}

// TODO: Implement return helper for inline bot result

// SendInlineBotResult invokes method messages.sendInlineBotResult#7aa11297 returning error if any. Send a result obtained using messages.getInlineBotResults¹.
func (ctx *Context) SendInlineBotResult(chatId int64, request tg.MessagesSendInlineBotResultRequest) (tg.UpdatesClass, error) {
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendInlineBotResult(ctx.OriginContext, &request)
}

// SendReaction invokes method messages.sendReaction#25690ce4 returning error if any.
func (ctx *Context) SendReaction(chatId int64, request tg.MessagesSendReactionRequest) (*tg.Message, error) {
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnNewMessageWithError(ctx.Client.MessagesSendReaction(ctx.OriginContext, &request))
}

// SendMultiMedia invokes method messages.sendMultiMedia#f803138f returning error if any. Send an album or grouped media¹
func (ctx *Context) SendMultiMedia(chatId int64, request tg.MessagesSendMultiMediaRequest) (*tg.Message, error) {
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnNewMessageWithError(ctx.Client.MessagesSendMultiMedia(ctx.OriginContext, &request))
}

// AnswerCallback invokes method messages.setBotCallbackAnswer#d58f130a returning error if any. Set the callback answer to a user button press
func (ctx *Context) AnswerCallback(request tg.MessagesSetBotCallbackAnswerRequest) (bool, error) {
	return ctx.Client.MessagesSetBotCallbackAnswer(ctx.OriginContext, &request)
}

// EditMessage invokes method messages.editMessage#48f71778 returning error if any. Edit message
func (ctx *Context) EditMessage(chatId int64, request tg.MessagesEditMessageRequest) (*tg.Message, error) {
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnEditMessageWithError(ctx.Client.MessagesEditMessage(ctx.OriginContext, &request))
}

// GetChat returns tg.ChatFullClass of the provided chat id.
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

// GetUser returns tg.UserFull of the provided user id.
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

// ExportSessionString returns session of authorized account in the form of string.
// Note: This session string can be used to log back in with the help of gotgproto.
// Check sessionMaker.SessionType for more information about it.
func (ctx *Context) ExportSessionString() (string, error) {
	return functions.EncodeSessionToString(storage.GetSession())
}

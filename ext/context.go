package ext

import (
	"context"
	"time"

	"github.com/anonyindian/gotgproto/functions"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// Context consists of context.Context, tg.Client, Self etc.
type Context struct {
	// original context of the client.
	context.Context
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
		Context:  ctx,
		Client:   client,
		Self:     self,
		Sender:   sender,
		Entities: entities,
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
	var m = &tg.Message{}
	switch text := (text).(type) {
	case string:
		m.Message = text
		u, err := builder.Text(ctx, text)
		return functions.ReturnNewMessageWithError(m, u, err)
	case []styling.StyledTextOption:
		tb := entity.Builder{}
		if err := styling.Perform(&tb, text...); err != nil {
			return nil, err
		}
		m.Message, _ = tb.Complete()
		u, err := builder.StyledText(ctx, text...)
		return functions.ReturnNewMessageWithError(m, u, err)
	default:
		return nil, ErrTextInvalid
	}
}

// SendMessage invokes method messages.sendMessage#d9d75a4 returning error if any.
func (ctx *Context) SendMessage(chatId int64, request *tg.MessagesSendMessageRequest) (*tg.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMessageRequest{}
	}
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	var m = &tg.Message{}
	m.Message = request.Message
	u, err := ctx.Client.MessagesSendMessage(ctx, request)
	return functions.ReturnNewMessageWithError(m, u, err)
}

// SendMedia invokes method messages.sendMedia#e25ff8e0 returning error if any. Send a media
func (ctx *Context) SendMedia(chatId int64, request *tg.MessagesSendMediaRequest) (*tg.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMediaRequest{}
	}
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	var m = &tg.Message{}
	m.Message = request.Message
	u, err := ctx.Client.MessagesSendMedia(ctx, request)
	return functions.ReturnNewMessageWithError(m, u, err)
}

// TODO: Implement return helper for inline bot result

// SendInlineBotResult invokes method messages.sendInlineBotResult#7aa11297 returning error if any. Send a result obtained using messages.getInlineBotResults¹.
func (ctx *Context) SendInlineBotResult(chatId int64, request *tg.MessagesSendInlineBotResultRequest) (tg.UpdatesClass, error) {
	if request == nil {
		request = &tg.MessagesSendInlineBotResultRequest{}
	}
	request.RandomID = time.Now().UnixNano()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return ctx.Client.MessagesSendInlineBotResult(ctx, request)
}

// SendReaction invokes method messages.sendReaction#25690ce4 returning error if any.
func (ctx *Context) SendReaction(chatId int64, request *tg.MessagesSendReactionRequest) (*tg.Message, error) {
	if request == nil {
		request = &tg.MessagesSendReactionRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	var m = &tg.Message{}
	m.Message = request.Reaction
	u, err := ctx.Client.MessagesSendReaction(ctx, request)
	return functions.ReturnNewMessageWithError(m, u, err)
}

// SendMultiMedia invokes method messages.sendMultiMedia#f803138f returning error if any. Send an album or grouped media¹
func (ctx *Context) SendMultiMedia(chatId int64, request *tg.MessagesSendMultiMediaRequest) (*tg.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMultiMediaRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	u, err := ctx.Client.MessagesSendMultiMedia(ctx, request)
	return functions.ReturnNewMessageWithError(&tg.Message{}, u, err)
}

// AnswerCallback invokes method messages.setBotCallbackAnswer#d58f130a returning error if any. Set the callback answer to a user button press
func (ctx *Context) AnswerCallback(request *tg.MessagesSetBotCallbackAnswerRequest) (bool, error) {
	if request == nil {
		request = &tg.MessagesSetBotCallbackAnswerRequest{}
	}
	return ctx.Client.MessagesSetBotCallbackAnswer(ctx, request)
}

// EditMessage invokes method messages.editMessage#48f71778 returning error if any. Edit message
func (ctx *Context) EditMessage(chatId int64, request *tg.MessagesEditMessageRequest) (*tg.Message, error) {
	if request == nil {
		request = &tg.MessagesEditMessageRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(chatId)
	}
	return functions.ReturnEditMessageWithError(ctx.Client.MessagesEditMessage(ctx, request))
}

// GetChat returns tg.ChatFullClass of the provided chat id.
func (ctx *Context) GetChat(chatId int64) (tg.ChatFullClass, error) {
	peer := storage.GetPeerById(chatId)
	if peer.ID == 0 {
		return nil, ErrPeerNotFound
	}
	switch storage.EntityType(peer.Type) {
	case storage.TypeChannel:
		channel, err := ctx.Client.ChannelsGetFullChannel(ctx, &tg.InputChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return channel.FullChat, nil
	case storage.TypeChat:
		chat, err := ctx.Client.MessagesGetFullChat(ctx, chatId)
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
	if peer.Type == storage.TypeUser.GetInt() {
		user, err := ctx.Client.UsersGetFullUser(ctx, &tg.InputUser{
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

func (ctx *Context) GetMessages(messageIds []tg.InputMessageClass) ([]tg.MessageClass, error) {
	return functions.GetMessages(ctx, ctx.Client, messageIds)
}

func (ctx *Context) BanChatMember(chatId, userId int64, untilDate int) (tg.UpdatesClass, error) {
	peerChatStorage := storage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return nil, ErrPeerNotFound
	}
	var chatPeer tg.InputPeerClass
	switch storage.EntityType(peerChatStorage.Type) {
	case storage.TypeChannel:
		chatPeer = &tg.InputPeerChannel{
			ChannelID:  peerChatStorage.ID,
			AccessHash: peerChatStorage.AccessHash,
		}
	case storage.TypeChat:
		chatPeer = &tg.InputPeerChat{
			ChatID: peerChatStorage.ID,
		}
	}
	peerUser := storage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return nil, ErrPeerNotFound
	}
	return functions.BanChatMember(ctx, ctx.Client, chatPeer, &tg.InputPeerUser{
		UserID:     peerUser.ID,
		AccessHash: peerUser.AccessHash,
	}, untilDate)
}

func (ctx *Context) UnbanChatMember(chatId, userId int64, untilDate int) (bool, error) {
	peerChatStorage := storage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return false, ErrPeerNotFound
	}
	var chatPeer = &tg.InputPeerChannel{}
	switch storage.EntityType(peerChatStorage.Type) {
	case storage.TypeChannel:
		chatPeer = &tg.InputPeerChannel{
			ChannelID:  peerChatStorage.ID,
			AccessHash: peerChatStorage.AccessHash,
		}
	default:
		return false, ErrNotChannel
	}
	peerUser := storage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return false, ErrPeerNotFound
	}
	return functions.UnbanChatMember(ctx, ctx.Client, chatPeer, &tg.InputPeerUser{
		UserID:     peerUser.ID,
		AccessHash: peerUser.AccessHash,
	})
}

func (ctx *Context) AddChatMembers(chatId int64, userIds []int64, forwardLimit int) (bool, error) {
	peerChatStorage := storage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return false, ErrPeerNotFound
	}
	var chatPeer tg.InputPeerClass
	switch storage.EntityType(peerChatStorage.Type) {
	case storage.TypeChannel:
		chatPeer = &tg.InputPeerChannel{
			ChannelID:  peerChatStorage.ID,
			AccessHash: peerChatStorage.AccessHash,
		}
	case storage.TypeChat:
		chatPeer = &tg.InputPeerChat{
			ChatID: peerChatStorage.ID,
		}
	default:
		return false, ErrNotChat
	}
	userPeers := make([]tg.InputUserClass, len(userIds))
	for i, uId := range userIds {
		userPeer := storage.GetPeerById(uId)
		if userPeer.ID == 0 {
			return false, ErrPeerNotFound
		}
		if userPeer.Type != int(storage.TypeUser) {
			return false, ErrNotUser
		}
		userPeers[i] = &tg.InputUser{
			UserID:     userPeer.ID,
			AccessHash: userPeer.AccessHash,
		}
	}
	return functions.AddChatMembers(ctx, ctx.Client, chatPeer, userPeers, forwardLimit)
}

func (ctx *Context) ArchiveChats(chatIds []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIds))
	for i, chatId := range chatIds {
		peer := storage.GetPeerById(chatId)
		if peer.ID == 0 {
			return false, ErrPeerNotFound
		}
		switch storage.EntityType(peer.Type) {
		case storage.TypeChannel:
			chatPeers[i] = &tg.InputPeerChannel{
				ChannelID:  peer.ID,
				AccessHash: peer.AccessHash,
			}
		case storage.TypeUser:
			chatPeers[i] = &tg.InputPeerUser{
				UserID:     peer.ID,
				AccessHash: peer.AccessHash,
			}
		case storage.TypeChat:
			chatPeers[i] = &tg.InputPeerChat{
				ChatID: peer.ID,
			}
		}
	}
	return functions.ArchiveChats(ctx, ctx.Client, chatPeers)
}

func (ctx *Context) UnarchiveChats(chatIds []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIds))
	for i, chatId := range chatIds {
		peer := storage.GetPeerById(chatId)
		if peer.ID == 0 {
			return false, ErrPeerNotFound
		}
		switch storage.EntityType(peer.Type) {
		case storage.TypeChannel:
			chatPeers[i] = &tg.InputPeerChannel{
				ChannelID:  peer.ID,
				AccessHash: peer.AccessHash,
			}
		case storage.TypeUser:
			chatPeers[i] = &tg.InputPeerUser{
				UserID:     peer.ID,
				AccessHash: peer.AccessHash,
			}
		case storage.TypeChat:
			chatPeers[i] = &tg.InputPeerChat{
				ChatID: peer.ID,
			}
		}
	}
	return functions.UnarchiveChats(ctx, ctx.Client, chatPeers)
}

func (ctx *Context) CreateChannel(title, about string, broadcast bool) (tg.UpdatesClass, error) {
	return functions.CreateChannel(ctx, ctx.Client, title, about, broadcast)
}

func (ctx *Context) CreateChat(title string, userIds []int64) (tg.UpdatesClass, error) {
	userPeers := make([]tg.InputUserClass, len(userIds))
	for i, uId := range userIds {
		userPeer := storage.GetPeerById(uId)
		if userPeer.ID == 0 {
			return nil, ErrPeerNotFound
		}
		if userPeer.Type != int(storage.TypeUser) {
			return nil, ErrNotUser
		}
		userPeers[i] = &tg.InputUser{
			UserID:     userPeer.ID,
			AccessHash: userPeer.AccessHash,
		}
	}
	return functions.CreateChat(ctx, ctx.Client, title, userPeers)
}

// TODO: Add documentation

func (ctx *Context) ResolveUsername(username string) (tg.PeerClass, error) {
	peer, err := ctx.Client.ContactsResolveUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return peer.Peer, nil
}

// ExportSessionString returns session of authorized account in the form of string.
// Note: This session string can be used to log back in with the help of gotgproto.
// Check sessionMaker.SessionType for more information about it.
func (ctx *Context) ExportSessionString() (string, error) {
	return functions.EncodeSessionToString(storage.GetSession())
}

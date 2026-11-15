package ext

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	mtp_errors "github.com/celestix/gotgproto/errors"
	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/storage"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// Context consists of context.Context, tg.Client, Self etc.
type Context struct {
	// raw tg client which will be used make send requests.
	Raw *tg.Client
	// self user who authorized the session.
	Self *tg.User
	// Sender is a message sending helper.
	Sender *message.Sender
	// Entities consists of mapped users, chats and channels from the update.
	Entities *tg.Entities
	// original context of the client.
	context.Context

	setReply    bool
	random      *rand.Rand
	PeerStorage *storage.PeerStorage
}

// NewContext creates a new Context object with provided parameters.
func NewContext(ctx context.Context, client *tg.Client, peerStorage *storage.PeerStorage, self *tg.User, sender *message.Sender, entities *tg.Entities, setReply bool) *Context {
	return &Context{
		Context:     ctx,
		Raw:         client,
		Self:        self,
		Sender:      sender,
		Entities:    entities,
		random:      rand.New(rand.NewSource(time.Now().UnixNano())),
		setReply:    setReply,
		PeerStorage: peerStorage,
	}
}

func (ctx *Context) generateRandomID() int64 {
	return ctx.random.Int63()
}

// ReplyOpts object contains optional parameters for Context.Reply.
type ReplyOpts struct {
	// Whether the message should show link preview or not.
	NoWebpage bool
	// Reply markup of a message, i.e. inline keyboard buttons etc.
	Markup           tg.ReplyMarkupClass
	ReplyToMessageId int
}

type ReplyTextType interface {
	construct()
}

type ReplyTextTypeString string

func (*ReplyTextTypeString) construct() {}

func (r ReplyTextTypeString) get() string {
	return string(r)
}

func ReplyTextString(s string) ReplyTextType {
	r := ReplyTextTypeString(s)
	return &r
}

type ReplyTextTypeStyledText styling.StyledTextOption

func (ReplyTextTypeStyledText) construct() {}

func (r ReplyTextTypeStyledText) get() styling.StyledTextOption {
	return styling.StyledTextOption(r)
}

func ReplyTextStyledText(s styling.StyledTextOption) ReplyTextType {
	r := ReplyTextTypeStyledText(s)
	return &r
}

type ReplyTextTypeStyledTextArray []styling.StyledTextOption

func (ReplyTextTypeStyledTextArray) construct() {}

func (r ReplyTextTypeStyledTextArray) get() []styling.StyledTextOption {
	return []styling.StyledTextOption(r)
}

func ReplyTextStyledTextArray(s []styling.StyledTextOption) ReplyTextType {
	r := ReplyTextTypeStyledTextArray(s)
	return &r
}

// Reply uses given message update to create message for same chat and create a reply.
// Parameter 'text' interface should be one from string or an array of styling.StyledTextOption.
func (ctx *Context) Reply(upd *Update, text ReplyTextType, opts *ReplyOpts) (*types.Message, error) {
	if text == nil {
		return nil, mtp_errors.ErrTextEmpty
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
	case *ReplyTextTypeString:
		m.Message = text.get()
		u, err := builder.Text(ctx, text.get())
		m, err = functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
		if err != nil {
			return nil, err
		}
	case *ReplyTextTypeStyledText:
		tb := entity.Builder{}
		if err := styling.Perform(&tb, text.get()); err != nil {
			return nil, err
		}
		m.Message, _ = tb.Complete()
		u, err := builder.StyledText(ctx, text.get())
		m, err = functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
		if err != nil {
			return nil, err
		}
	case *ReplyTextTypeStyledTextArray:
		tb := entity.Builder{}
		if err := styling.Perform(&tb, text.get()...); err != nil {
			return nil, err
		}
		m.Message, _ = tb.Complete()
		u, err := builder.StyledText(ctx, text.get()...)
		m, err = functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
		if err != nil {
			return nil, err
		}
	default:
		return nil, mtp_errors.ErrTextInvalid
	}
	msg := types.ConstructMessage(m)
	msg.ReplyToMessage = upd.EffectiveMessage
	return msg, nil
}

// SendMessage invokes method messages.sendMessage#d9d75a4 returning error if any.
func (ctx *Context) SendMessage(chatID int64, request *tg.MessagesSendMessageRequest) (*types.Message, error) {
	var err error
	if request == nil {
		request = &tg.MessagesSendMessageRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer, err = ctx.ResolveInputPeerById(chatID)
		if err != nil {
			return nil, err
		}
	}
	var m = &tg.Message{}
	m.Message = request.Message
	u, err := ctx.Raw.MessagesSendMessage(ctx, request)
	message, err := functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
	if err != nil {
		return nil, err
	}
	msg := types.ConstructMessage(message)
	if ctx.setReply {
		_ = msg.SetRepliedToMessage(ctx.Context, ctx.Raw, ctx.PeerStorage)
	}
	return msg, nil
}

// SendMedia invokes method messages.sendMedia#e25ff8e0 returning error if any. Send a media
func (ctx *Context) SendMedia(chatID int64, request *tg.MessagesSendMediaRequest) (*types.Message, error) {
	var err error
	if request == nil {
		request = &tg.MessagesSendMediaRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer, err = ctx.ResolveInputPeerById(chatID)
		if err != nil {
			return nil, err
		}
	}
	var m = &tg.Message{}
	m.Message = request.Message
	u, err := ctx.Raw.MessagesSendMedia(ctx, request)
	message, err := functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
	if err != nil {
		return nil, err
	}
	msg := types.ConstructMessage(message)
	if ctx.setReply {
		_ = msg.SetRepliedToMessage(ctx.Context, ctx.Raw, ctx.PeerStorage)
	}
	return msg, nil
}

// SetInlineBotResult invokes method messages.setInlineBotResults#eb5ea206 returning error if any.
// Answer an inline query, for bots only
func (ctx *Context) SetInlineBotResult(request *tg.MessagesSetInlineBotResultsRequest) (bool, error) {
	return ctx.Raw.MessagesSetInlineBotResults(ctx, request)
}

func (ctx *Context) GetInlineBotResults(chatID int64, botUsername string, request *tg.MessagesGetInlineBotResultsRequest) (*tg.MessagesBotResults, error) {
	bot := ctx.PeerStorage.GetPeerByUsername(botUsername)
	if bot.ID == 0 {
		c, err := ctx.ResolveUsername(botUsername)
		if err != nil {
			return nil, err
		}
		switch {
		case c.IsAUser():
			bot = &storage.Peer{
				ID:         c.GetID(),
				AccessHash: c.GetAccessHash(),
			}
		default:
			return nil, errors.New("provided username was invalid for a bot")
		}
	}
	request.Peer, _ = ctx.ResolveInputPeerById(chatID)
	request.Bot = &tg.InputUser{
		UserID:     bot.ID,
		AccessHash: bot.AccessHash,
	}
	return ctx.Raw.MessagesGetInlineBotResults(ctx, request)
}

// TODO: Implement return helper for inline bot result

// SendInlineBotResult invokes method messages.sendInlineBotResult#7aa11297 returning error if any. Send a result obtained using messages.getInlineBotResults¹.
func (ctx *Context) SendInlineBotResult(chatID int64, request *tg.MessagesSendInlineBotResultRequest) (tg.UpdatesClass, error) {
	if request == nil {
		request = &tg.MessagesSendInlineBotResultRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer, _ = ctx.ResolveInputPeerById(chatID)
	}
	return ctx.Raw.MessagesSendInlineBotResult(ctx, request)
}

// SendReaction invokes messages.sendReaction. It returns a nil message when the
// successful response does not contain a new message.
func (ctx *Context) SendReaction(chatID int64, request *tg.MessagesSendReactionRequest) (*types.Message, error) {
	var err error
	if request == nil {
		request = &tg.MessagesSendReactionRequest{}
	}
	if request.Peer == nil {
		request.Peer, err = ctx.ResolveInputPeerById(chatID)
		if err != nil {
			return nil, err
		}
	}
	var m = &tg.Message{}
	// m.Message = request.Reaction
	u, err := ctx.Raw.MessagesSendReaction(ctx, request)
	message, err := functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, nil
	}
	msg := types.ConstructMessage(message)
	if ctx.setReply {
		_ = msg.SetRepliedToMessage(ctx.Context, ctx.Raw, ctx.PeerStorage)
	}
	return msg, nil
}

// SendMultiMedia invokes method messages.sendMultiMedia#f803138f returning error if any. Send an album or grouped media¹
func (ctx *Context) SendMultiMedia(chatID int64, request *tg.MessagesSendMultiMediaRequest) (*types.Message, error) {
	var err error
	if request == nil {
		request = &tg.MessagesSendMultiMediaRequest{}
	}
	if request.Peer == nil {
		request.Peer, err = ctx.ResolveInputPeerById(chatID)
		if err != nil {
			return nil, err
		}
	}
	u, err := ctx.Raw.MessagesSendMultiMedia(ctx, request)
	message, err := functions.ReturnNewMessageWithError(&tg.Message{}, u, ctx.PeerStorage, err)
	if err != nil {
		return nil, err
	}
	msg := types.ConstructMessage(message)
	if ctx.setReply {
		_ = msg.SetRepliedToMessage(ctx.Context, ctx.Raw, ctx.PeerStorage)
	}
	return msg, nil
}

// AnswerCallback invokes method messages.setBotCallbackAnswer#d58f130a returning error if any. Set the callback answer to a user button press
func (ctx *Context) AnswerCallback(request *tg.MessagesSetBotCallbackAnswerRequest) (bool, error) {
	if request == nil {
		request = &tg.MessagesSetBotCallbackAnswerRequest{}
	}
	return ctx.Raw.MessagesSetBotCallbackAnswer(ctx, request)
}

// EditMessage invokes method messages.editMessage#48f71778 returning error if any. Edit message
func (ctx *Context) EditMessage(chatID int64, request *tg.MessagesEditMessageRequest) (*types.Message, error) {
	var err error
	if request == nil {
		request = &tg.MessagesEditMessageRequest{}
	}
	if request.Peer == nil {
		request.Peer, err = ctx.ResolveInputPeerById(chatID)
		if err != nil {
			return nil, err
		}
	}
	upds, err := ctx.Raw.MessagesEditMessage(ctx, request)
	message, err := functions.ReturnEditMessageWithError(ctx.PeerStorage, upds, err)
	if err != nil {
		return nil, err
	}
	msg := types.ConstructMessage(message)
	if ctx.setReply {
		_ = msg.SetRepliedToMessage(ctx.Context, ctx.Raw, ctx.PeerStorage)
	}
	return msg, nil
}

// GetChat returns tg.ChatFullClass of the provided chat id.
func (ctx *Context) GetChat(chatID int64) (tg.ChatFullClass, error) {
	inputPeer, err := ctx.ResolveInputPeerById(chatID)
	if err != nil  {
		return nil, err
	}
	switch p := inputPeer.(type) {
    case *tg.InputPeerChannel:
        channel, err := ctx.Raw.ChannelsGetFullChannel(ctx, &tg.InputChannel{
            ChannelID:  p.ChannelID,
            AccessHash: p.AccessHash,
        })
        if err != nil {
            return nil, err
        }
        return channel.FullChat, nil
    case *tg.InputPeerChat:
        chat, err := ctx.Raw.MessagesGetFullChat(ctx, chatID)
        if err != nil {
            return nil, err
        }
        return chat.FullChat, nil
    case *tg.InputPeerEmpty:
        return nil, mtp_errors.ErrPeerNotFound
    default:
        return nil, mtp_errors.ErrNotChat
    }
}

// GetUser returns tg.UserFull of the provided user id.
func (ctx *Context) GetUser(userID int64) (*tg.UserFull, error) {
	inputPeer, err := ctx.ResolveInputPeerById(userID)
	if err != nil  {
		return nil, err
	}
	switch p := inputPeer.(type) {
	case *tg.InputPeerUser:
		user, err := ctx.Raw.UsersGetFullUser(ctx, &tg.InputUser{
			UserID:     p.UserID,
			AccessHash: p.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return &user.FullUser, nil
	default:
		return nil, mtp_errors.ErrNotUser
	}
}

// GetMessages is used to fetch messages from a PM (Private Chat).
func (ctx *Context) GetMessages(chatID int64, messageIDs []tg.InputMessageClass) ([]tg.MessageClass, error) {
	return functions.GetMessages(ctx.Context, ctx.Raw, ctx.PeerStorage, chatID, messageIDs)
}

// BanChatMember is used to ban a user from a chat.
func (ctx *Context) BanChatMember(chatID, userID int64, untilDate int) (tg.UpdatesClass, error) {
	inputPeerChat, err :=  ctx.ResolveInputPeerById(chatID)
	if err != nil  {
		return nil, err
	}
	switch inputPeerChat.(type) {
    case *tg.InputPeerChannel:
    case *tg.InputPeerChat:
    case *tg.InputPeerEmpty:
			return nil, mtp_errors.ErrPeerNotFound
    default:
			return nil, mtp_errors.ErrNotChat
    }
	var inputPeerUser *tg.InputPeerUser
	inputPeer, err := ctx.ResolveInputPeerById(userID)
	if err != nil  {
		return nil, err
	}
	switch p := inputPeer.(type) {
	case *tg.InputPeerUser:
		inputPeerUser = p
	case *tg.InputPeerEmpty:
		return nil, mtp_errors.ErrPeerNotFound
	default:
		return nil, mtp_errors.ErrNotUser
	}
	return functions.BanChatMember(ctx, ctx.Raw, inputPeerChat, inputPeerUser, untilDate)
}

// UnbanChatMember is used to unban a user from a chat.
func (ctx *Context) UnbanChatMember(chatID, userID int64) (bool, error) {
	var inputPeerChat *tg.InputPeerChannel
	inputPeer, err :=  ctx.ResolveInputPeerById(chatID)
	if err != nil  {
		return false, err
	}
	switch p := inputPeer.(type) {
    case *tg.InputPeerChannel:
			inputPeerChat = p
    case *tg.InputPeerEmpty:
			return false, mtp_errors.ErrPeerNotFound
    default:
			return false, mtp_errors.ErrNotChannel
    }
	var inputPeerUser *tg.InputPeerUser
	inputPeer, err = ctx.ResolveInputPeerById(userID)
	if err != nil  {
		return false, err
	}
	switch p := inputPeer.(type) {
	
	case *tg.InputPeerUser:
		inputPeerUser= p
	case *tg.InputPeerEmpty:
    return false, mtp_errors.ErrPeerNotFound
	default:
		return false, mtp_errors.ErrNotUser
	}
	return functions.UnbanChatMember(ctx, ctx.Raw, inputPeerChat, inputPeerUser)
}

// AddChatMembers is used to add members to a chat
func (ctx *Context) AddChatMembers(chatID int64, userIDs []int64, forwardLimit int) (bool, error) {
	inputPeerChat, err :=  ctx.ResolveInputPeerById(chatID)
	if err != nil  {
		return false, err
	}
	switch inputPeerChat.(type) {
    case *tg.InputPeerChannel:
    case *tg.InputPeerChat:
    case *tg.InputPeerEmpty:
        return false, mtp_errors.ErrPeerNotFound
    default:
        return false, mtp_errors.ErrNotChat
    }
	userPeers := make([]tg.InputUserClass, len(userIDs))
	for i, uID := range userIDs {
		inputPeerUser, err := ctx.ResolveInputPeerById(uID)
		if err != nil  {
			return false, err
		}
		switch p := inputPeerUser.(type) {
		case *tg.InputPeerUser:
			userPeers[i] = &tg.InputUser{
				UserID:     p.UserID,
				AccessHash: p.AccessHash,
			}
		case *tg.InputPeerEmpty:
	        return false, mtp_errors.ErrPeerNotFound
		default:
			return false, mtp_errors.ErrNotUser
		}
	}
	return functions.AddChatMembers(ctx, ctx.Raw, inputPeerChat, userPeers, forwardLimit)
}

// ArchiveChats invokes method folders.editPeerFolders#6847d0ab returning error if any.
// Edit peers in peer folder¹
//
// Links:
//  1. https://core.telegram.org/api/folders#peer-folders
func (ctx *Context) ArchiveChats(chatIDs []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIDs))
	for i, chatID := range chatIDs {
		inputPeer, err := ctx.ResolveInputPeerById(chatID)
		if err != nil  {
			return false, err
		}
		switch inputPeer.(type) {
	    case *tg.InputPeerChannel:
	    case *tg.InputPeerChat:
	    case *tg.InputPeerUser:
	    case *tg.InputPeerEmpty:
	        return false, mtp_errors.ErrPeerNotFound
	    default:
	        return false, mtp_errors.ErrNotChat
	    }
		chatPeers[i] = inputPeer
	}
	return functions.ArchiveChats(ctx, ctx.Raw, chatPeers)
}

// UnarchiveChats invokes method folders.editPeerFolders#6847d0ab returning error if any.
// Edit peers in peer folder¹
//
// Links:
//  1. https://core.telegram.org/api/folders#peer-folders
func (ctx *Context) UnarchiveChats(chatIDs []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIDs))
	for i, chatID := range chatIDs {
		inputPeer, err := ctx.ResolveInputPeerById(chatID)
		if err != nil  {
			return false, err
		}
		switch inputPeer.(type) {
	    case *tg.InputPeerChannel:
	    case *tg.InputPeerChat:
	    case *tg.InputPeerUser:
	    case *tg.InputPeerEmpty:
	        return false, mtp_errors.ErrPeerNotFound
	    default:
	        return false, mtp_errors.ErrNotChat
	    }
		chatPeers[i] = inputPeer
	}
	return functions.UnarchiveChats(ctx, ctx.Raw, chatPeers)
}

// CreateChannel invokes method channels.createChannel#3d5fb10f returning error if any.
// Create a supergroup/channel¹.
//
// Links:
//  1. https://core.telegram.org/api/channel
func (ctx *Context) CreateChannel(title, about string, broadcast bool) (*tg.Channel, error) {
	return functions.CreateChannel(ctx, ctx.Raw, ctx.PeerStorage, title, about, broadcast)
}

// CreateChat invokes method messages.createChat#9cb126e returning error if any. Creates a new chat.
func (ctx *Context) CreateChat(title string, userIDs []int64) (*tg.Chat, error) {
	userPeers := make([]tg.InputUserClass, len(userIDs))
	for i, uID := range userIDs {
		userPeer := ctx.ResolvePeerById(uID)
		if userPeer.ID == 0 {
			return nil, mtp_errors.ErrPeerNotFound
		}
		if userPeer.Type != int(storage.TypeUser) {
			return nil, mtp_errors.ErrNotUser
		}
		userPeers[i] = &tg.InputUser{
			UserID:     userPeer.ID,
			AccessHash: userPeer.AccessHash,
		}
	}
	return functions.CreateChat(ctx, ctx.Raw, ctx.PeerStorage, title, userPeers)
}

// DeleteMessages shall be used to delete messages in a chat with chatID and messageIDs.
// Returns error if failed to delete.
func (ctx *Context) DeleteMessages(chatID int64, messageIDs []int) error {
	inputPeer, err := ctx.ResolveInputPeerById(chatID)
	if err != nil  {
		return err
	}
	switch p := inputPeer.(type) {
    case *tg.InputPeerChannel:
        _, err := ctx.Raw.ChannelsDeleteMessages(ctx, &tg.ChannelsDeleteMessagesRequest{
			Channel: &tg.InputChannel{
				ChannelID:  p.ChannelID,
				AccessHash: p.AccessHash,
			},
			ID: messageIDs,
		})
		return err
    case *tg.InputPeerChat, *tg.InputPeerUser:
        _, err := ctx.Raw.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     messageIDs,
		})
		return err
    case *tg.InputPeerEmpty:
        return mtp_errors.ErrPeerNotFound
    default:
        return mtp_errors.ErrNotChat
    }
}

// ForwardMessage shall be used to forward messages in a chat with chatID and messageIDs.
// Returns updatesclass or an error if failed to delete.
//
// Deprecated: use ForwardMessages instead.
func (ctx *Context) ForwardMessage(fromChatID, toChatID int64, request *tg.MessagesForwardMessagesRequest) (tg.UpdatesClass, error) {
	return ctx.ForwardMessages(fromChatID, toChatID, request)
}

// ForwardMessages shall be used to forward messages in a chat with chatID and messageIDs.
// Returns updatesclass or an error if failed to delete.
func (ctx *Context) ForwardMessages(fromChatID, toChatID int64, request *tg.MessagesForwardMessagesRequest) (tg.UpdatesClass, error) {
	fromPeer, _ := ctx.ResolveInputPeerById(fromChatID)
	if fromPeer.Zero() {
		return nil, fmt.Errorf("fromChatID: %w", mtp_errors.ErrPeerNotFound)
	}
	toPeer, _ := ctx.ResolveInputPeerById(toChatID)
	if toPeer.Zero() {
		return nil, fmt.Errorf("toChatID: %w", mtp_errors.ErrPeerNotFound)
	}
	if request == nil {
		request = &tg.MessagesForwardMessagesRequest{}
	}
	request.RandomID = make([]int64, len(request.ID))
	for i := 0; i < len(request.ID); i++ {
		request.RandomID[i] = ctx.generateRandomID()
	}
	return ctx.Raw.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{
		RandomID:           request.RandomID,
		ID:                 request.ID,
		FromPeer:           fromPeer,
		ToPeer:             toPeer,
		DropAuthor:         request.DropAuthor,
		Silent:             request.Silent,
		Background:         request.Background,
		WithMyScore:        request.WithMyScore,
		DropMediaCaptions:  request.DropMediaCaptions,
		Noforwards:         request.Noforwards,
		TopMsgID:           request.TopMsgID,
		ScheduleDate:       request.ScheduleDate,
		SendAs:             request.SendAs,
		QuickReplyShortcut: request.QuickReplyShortcut,
	})
}

type EditAdminOpts struct {
	AdminRights tg.ChatAdminRights
	AdminTitle  string
}

// PromoteChatMember is used to promote a user in a chat.
func (ctx *Context) PromoteChatMember(chatID, userID int64, opts *EditAdminOpts) (bool, error) {
	peerChat := ctx.ResolvePeerById(chatID)
	if peerChat.ID == 0 {
		return false, fmt.Errorf("chat: %w", mtp_errors.ErrPeerNotFound)
	}
	peerUser := ctx.ResolvePeerById(userID)
	if peerUser.ID == 0 {
		return false, fmt.Errorf("user: %w", mtp_errors.ErrPeerNotFound)
	}
	if opts == nil {
		opts = &EditAdminOpts{}
	}
	return functions.PromoteChatMember(ctx, ctx.Raw, peerChat, peerUser, opts.AdminRights, opts.AdminTitle)
}

// DemoteChatMember is used to demote a user in a chat.
func (ctx *Context) DemoteChatMember(chatID, userID int64, opts *EditAdminOpts) (bool, error) {
	peerChat := ctx.ResolvePeerById(chatID)
	if peerChat.ID == 0 {
		return false, fmt.Errorf("chat: %w", mtp_errors.ErrPeerNotFound)
	}
	peerUser := ctx.ResolvePeerById(userID)
	if peerUser.ID == 0 {
		return false, fmt.Errorf("user: %w", mtp_errors.ErrPeerNotFound)
	}
	if opts == nil {
		opts = &EditAdminOpts{}
	}
	return functions.DemoteChatMember(ctx, ctx.Raw, peerChat, peerUser, opts.AdminRights, opts.AdminTitle)
}

// ResolveUsername invokes method contacts.resolveUsername#f93ccba3 returning error if any.
// Resolve a @username to get peer info
func (ctx *Context) ResolveUsername(username string) (types.EffectiveChat, error) {
	return ctx.extractContactResolvedPeer(
		ctx.Raw.ContactsResolveUsername(
			ctx,
			&tg.ContactsResolveUsernameRequest{
				Username: strings.TrimPrefix(
					username,
					"@",
				),
			},
		),
	)
}

func (ctx *Context) extractContactResolvedPeer(p *tg.ContactsResolvedPeer, err error) (types.EffectiveChat, error) {
	if err != nil {
		return &types.EmptyUC{}, err
	}
	functions.SavePeersFromClassArray(ctx.PeerStorage, p.Chats, p.Users)
	switch p.Peer.(type) {
	case *tg.PeerChannel:
		if len(p.Chats) == 0 {
			return &types.EmptyUC{}, errors.New("peer info not found in the resolved Chats")
		}
		switch chat := p.Chats[0].(type) {
		case *tg.Channel:
			var c = types.Channel(*chat)
			return &c, nil
		case *tg.ChannelForbidden:
			return &types.EmptyUC{}, errors.New("peer could not be resolved because Channel Forbidden")
		}
	case *tg.PeerUser:
		if len(p.Users) == 0 {
			return &types.EmptyUC{}, errors.New("peer info not found in the resolved Chats")
		}
		switch user := p.Users[0].(type) {
		case *tg.User:
			var c = types.User(*user)
			return &c, nil
		}
	}
	return &types.EmptyUC{}, errors.New("contact not found")
}

// GetUserProfilePhotos invokes method photos.getUserPhotos#91cd32a8 returning error if any. Returns the list of user photos.
func (ctx *Context) GetUserProfilePhotos(userID int64, opts *tg.PhotosGetUserPhotosRequest) ([]tg.PhotoClass, error) {
	peerUser := ctx.ResolvePeerById(userID)
	if peerUser.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
	}
	if opts == nil {
		opts = &tg.PhotosGetUserPhotosRequest{}
	}
	opts.UserID = &tg.InputUser{
		UserID:     userID,
		AccessHash: peerUser.AccessHash,
	}
	p, err := ctx.Raw.PhotosGetUserPhotos(ctx, opts)
	if err != nil {
		return nil, err
	}
	return p.GetPhotos(), nil
}

// ExportSessionString returns session of authorized account in the form of string.
// Note: This session string can be used to log back in with the help of gotgproto.
// Check sessionMaker.SessionType for more information about it.
func (ctx *Context) ExportSessionString() (string, error) {
	return functions.EncodeSessionToString(ctx.PeerStorage.GetSession())
}

// DownloadOutputClass is an interface which is used to download media.
// It can be one from DownloadOutputStream, DownloadOutputPath and DownloadOutputParallel.
type DownloadOutputClass interface {
	run(context.Context, *downloader.Builder) (tg.StorageFileTypeClass, error)
}

// DownloadOutputStream is used to download media to an io.Writer.
// It can be used to download media to a file, memory etc.
type DownloadOutputStream struct {
	io.Writer
}

func (d DownloadOutputStream) run(ctx context.Context, b *downloader.Builder) (tg.StorageFileTypeClass, error) {
	return b.Stream(ctx, d)
}

// DownloadOutputPath is used to download media to a file.
type DownloadOutputPath string

func (d DownloadOutputPath) run(ctx context.Context, b *downloader.Builder) (tg.StorageFileTypeClass, error) {
	return b.ToPath(ctx, string(d))
}

// DownloadOutputParallel is used to download a media parallely.
type DownloadOutputParallel struct {
	io.WriterAt
}

func (d DownloadOutputParallel) run(ctx context.Context, b *downloader.Builder) (tg.StorageFileTypeClass, error) {
	return b.Parallel(ctx, d)
}

// DownloadMediaOpts object contains optional parameters for Context.DownloadMedia.
// If not provided, default values will be used.
type DownloadMediaOpts struct {
	// Threads sets downloading goroutines limit.
	Threads int
	//  If verify is true, file hashes will be checked. Verify is true by default for CDN downloads.
	Verify *bool
	// PartSize sets chunk size. Must be divisible by 4KB.
	PartSize int
}

// DownloadMedia downloads media from the provided MessageMediaClass.
// DownloadOutputClass can be one from DownloadOutputStream, DownloadOutputPath and DownloadOutputParallel.
// DownloadMediaOpts can be used to set optional parameters.
// Returns tg.StorageFileTypeClass and error if any.
func (ctx *Context) DownloadMedia(media tg.MessageMediaClass, downloadOutput DownloadOutputClass, opts *DownloadMediaOpts) (tg.StorageFileTypeClass, error) {
	if opts == nil {
		opts = &DownloadMediaOpts{}
	}
	mediaDownloader := downloader.NewDownloader()
	if opts.PartSize > 0 {
		mediaDownloader.WithPartSize(opts.PartSize)
	}
	inputFileLocation, err := functions.GetInputFileLocation(media)
	if err != nil {
		return nil, err
	}
	d := mediaDownloader.Download(ctx.Raw, inputFileLocation)
	if opts.Threads > 0 {
		d.WithThreads(opts.Threads)
	}
	if opts.Verify != nil {
		d.WithVerify(*opts.Verify)
	}
	return downloadOutput.run(ctx, d)
}

// TransferStarGift is used to transfer a star gift to a chat.
// Returns tg.UpdatesClass and error if any.
func (ctx *Context) TransferStarGift(chatID int64, starGift tg.InputSavedStarGiftClass) (tg.UpdatesClass, error) {
	peerUser, err := ctx.ResolveInputPeerById(chatID)
	if err != nil {
		return nil, mtp_errors.ErrPeerNotFound
	}
	upd, err := ctx.Raw.PaymentsTransferStarGift(ctx, &tg.PaymentsTransferStarGiftRequest{
		ToID:     peerUser,
		Stargift: starGift,
	})
	if err != nil {
		return nil, err
	}
	return upd, err
}

func (ctx *Context) ExportInvoice(inputMedia tg.InputMediaClass) (*tg.PaymentsExportedInvoice, error) {
	return ctx.Raw.PaymentsExportInvoice(ctx, inputMedia)
}

func (ctx *Context) SetPreCheckoutResults(success bool, queryID int64, err string) (bool, error) {
	return ctx.Raw.MessagesSetBotPrecheckoutResults(ctx, &tg.MessagesSetBotPrecheckoutResultsRequest{
		Success: success,
		QueryID: queryID,
		Error:   err,
	})
}

// ResolveInputPeerById tries to resolve given id to InputPeer.
// Returns tg.InputPeerClass or error if peer could not be resolved 
func (ctx *Context) ResolveInputPeerById(id int64) (tg.InputPeerClass, error) {
	peerStorage := ctx.PeerStorage
	peer := peerStorage.GetInputPeerById(id)
	if _, isEmpty := peer.(*tg.InputPeerEmpty); !isEmpty {
		return peer, nil
	}
	
	ID := constant.TDLibPeerID(id)
	if ID.IsChannel() {
		plainID := ID.ToPlain()
		chatsClass, err := ctx.Raw.ChannelsGetChannels(ctx, []tg.InputChannelClass{
			&tg.InputChannel{
				ChannelID:  plainID,
				AccessHash: 0,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch channel: %w", err)
		}
		chat, ok := chatsClass.MapChats().First()
		if ok {
			if ch, ok := chat.(*tg.Channel); ok {
				peerStorage.AddPeer(plainID, ch.AccessHash, storage.TypeChannel, ch.Username)
				return ch.AsInputPeer(), nil
			}
		}
	} else if ID.IsChat() {
		plainID := ID.ToPlain()
		peer := &tg.InputPeerChat{
			ChatID: plainID,
		}
		peerStorage.AddPeer(plainID, storage.DefaultAccessHash, storage.TypeChat, storage.DefaultUsername)
		return peer, nil
	} else {
		plainID := ID.ToPlain()
		users, err := ctx.Raw.UsersGetUsers(ctx, []tg.InputUserClass{
			&tg.InputUser{
				UserID:     plainID,
				AccessHash: 0,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user: %w", err)
		}
		user, ok := tg.UserClassArray(users).FirstAsNotEmpty()
		if ok {
			peerStorage.AddPeer(plainID, user.AccessHash, storage.TypeUser, user.Username)
			return user.AsInputPeer(), nil
		}
		//Try to get from storage again, but this time with bot-api compatible ids
		if ID.IsUser() {
			ID.Channel(id)
			peer := peerStorage.GetInputPeerById(int64(ID))
			if _, isEmpty := peer.(*tg.InputPeerEmpty); !isEmpty {
				return peer, nil
			}
			ID.Chat(id)
			peer = peerStorage.GetInputPeerById(int64(ID))
			if _, isEmpty := peer.(*tg.InputPeerEmpty); !isEmpty {
				return peer, nil
			}
		}
	}

	return nil, mtp_errors.ErrPeerNotFound
}


// ResolvePeerById tries to resolve given id to peer.
func (ctx *Context) ResolvePeerById(id int64) *storage.Peer {
	_, _ = ctx.ResolveInputPeerById(id)
	peer := ctx.PeerStorage.GetPeerById(id)
	if peer.ID != 0 {
		return peer
	}
	ID := constant.TDLibPeerID(id)
	if ID.IsUser() {
		ID.Channel(id)
		peer = ctx.ResolvePeerById(int64(ID))
		if peer.ID != 0 {
			return peer
		}
		ID.Chat(id)
		peer = ctx.ResolvePeerById(int64(ID))
		if peer.ID != 0 {
			return peer
		}
	
	}
	return peer
}

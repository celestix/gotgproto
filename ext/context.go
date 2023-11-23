package ext

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	mtp_errors "github.com/celestix/gotgproto/errors"
	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/storage"
	"github.com/celestix/gotgproto/types"
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
		random:      rand.New(rand.NewSource(time.Now().Unix())),
		setReply:    setReply,
		PeerStorage: peerStorage,
	}
}

func (c *Context) generateRandomID() int64 {
	return c.random.Int63()
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
func (ctx *Context) Reply(upd *Update, text interface{}, opts *ReplyOpts) (*types.Message, error) {
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
	case string:
		m.Message = text
		u, err := builder.Text(ctx, text)
		m, err = functions.ReturnNewMessageWithError(m, u, ctx.PeerStorage, err)
		if err != nil {
			return nil, err
		}
	case []styling.StyledTextOption:
		tb := entity.Builder{}
		if err := styling.Perform(&tb, text...); err != nil {
			return nil, err
		}
		m.Message, _ = tb.Complete()
		u, err := builder.StyledText(ctx, text...)
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
func (ctx *Context) SendMessage(chatId int64, request *tg.MessagesSendMessageRequest) (*types.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMessageRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
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
func (ctx *Context) SendMedia(chatId int64, request *tg.MessagesSendMediaRequest) (*types.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMediaRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
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

func (ctx *Context) GetInlineBotResults(chatId int64, botUsername string, request *tg.MessagesGetInlineBotResultsRequest) (*tg.MessagesBotResults, error) {
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
	request.Peer = ctx.PeerStorage.GetInputPeerById(chatId)
	request.Bot = &tg.InputUser{
		UserID:     bot.ID,
		AccessHash: bot.AccessHash,
	}
	return ctx.Raw.MessagesGetInlineBotResults(ctx, request)
}

// TODO: Implement return helper for inline bot result

// SendInlineBotResult invokes method messages.sendInlineBotResult#7aa11297 returning error if any. Send a result obtained using messages.getInlineBotResults¹.
func (ctx *Context) SendInlineBotResult(chatId int64, request *tg.MessagesSendInlineBotResultRequest) (tg.UpdatesClass, error) {
	if request == nil {
		request = &tg.MessagesSendInlineBotResultRequest{}
	}
	request.RandomID = ctx.generateRandomID()
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
	}
	return ctx.Raw.MessagesSendInlineBotResult(ctx, request)
}

// SendReaction invokes method messages.sendReaction#25690ce4 returning error if any.
func (ctx *Context) SendReaction(chatId int64, request *tg.MessagesSendReactionRequest) (*types.Message, error) {
	if request == nil {
		request = &tg.MessagesSendReactionRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
	}
	var m = &tg.Message{}
	// m.Message = request.Reaction
	u, err := ctx.Raw.MessagesSendReaction(ctx, request)
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

// SendMultiMedia invokes method messages.sendMultiMedia#f803138f returning error if any. Send an album or grouped media¹
func (ctx *Context) SendMultiMedia(chatId int64, request *tg.MessagesSendMultiMediaRequest) (*types.Message, error) {
	if request == nil {
		request = &tg.MessagesSendMultiMediaRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
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
func (ctx *Context) EditMessage(chatId int64, request *tg.MessagesEditMessageRequest) (*types.Message, error) {
	if request == nil {
		request = &tg.MessagesEditMessageRequest{}
	}
	if request.Peer == nil {
		request.Peer = functions.GetInputPeerClassFromId(ctx.PeerStorage, chatId)
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
func (ctx *Context) GetChat(chatId int64) (tg.ChatFullClass, error) {
	peer := ctx.PeerStorage.GetPeerById(chatId)
	if peer.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
	}
	switch storage.EntityType(peer.Type) {
	case storage.TypeChannel:
		channel, err := ctx.Raw.ChannelsGetFullChannel(ctx, &tg.InputChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return channel.FullChat, nil
	case storage.TypeChat:
		chat, err := ctx.Raw.MessagesGetFullChat(ctx, chatId)
		if err != nil {
			return nil, err
		}
		return chat.FullChat, nil
	}
	return nil, mtp_errors.ErrNotChat
}

// GetUser returns tg.UserFull of the provided user id.
func (ctx *Context) GetUser(userId int64) (*tg.UserFull, error) {
	peer := ctx.PeerStorage.GetPeerById(userId)
	if peer.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
	}
	if peer.Type == storage.TypeUser.GetInt() {
		user, err := ctx.Raw.UsersGetFullUser(ctx, &tg.InputUser{
			UserID:     peer.ID,
			AccessHash: peer.AccessHash,
		})
		if err != nil {
			return nil, err
		}
		return &user.FullUser, nil
	} else {
		return nil, mtp_errors.ErrNotUser
	}
}

// GetMessages is used to fetch messages from a PM (Private Chat).
func (ctx *Context) GetMessages(chatId int64, messageIds []tg.InputMessageClass) ([]tg.MessageClass, error) {
	return functions.GetMessages(ctx.Context, ctx.Raw, ctx.PeerStorage, chatId, messageIds)
}

// BanChatMember is used to ban a user from a chat.
func (ctx *Context) BanChatMember(chatId, userId int64, untilDate int) (tg.UpdatesClass, error) {
	peerChatStorage := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
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
	peerUser := ctx.PeerStorage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
	}
	return functions.BanChatMember(ctx, ctx.Raw, chatPeer, &tg.InputPeerUser{
		UserID:     peerUser.ID,
		AccessHash: peerUser.AccessHash,
	}, untilDate)
}

// UnbanChatMember is used to unban a user from a chat.
func (ctx *Context) UnbanChatMember(chatId, userId int64) (bool, error) {
	peerChatStorage := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return false, mtp_errors.ErrPeerNotFound
	}
	var chatPeer *tg.InputPeerChannel
	switch storage.EntityType(peerChatStorage.Type) {
	case storage.TypeChannel:
		chatPeer = &tg.InputPeerChannel{
			ChannelID:  peerChatStorage.ID,
			AccessHash: peerChatStorage.AccessHash,
		}
	default:
		return false, mtp_errors.ErrNotChannel
	}
	peerUser := ctx.PeerStorage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return false, mtp_errors.ErrPeerNotFound
	}
	return functions.UnbanChatMember(ctx, ctx.Raw, chatPeer, &tg.InputPeerUser{
		UserID:     peerUser.ID,
		AccessHash: peerUser.AccessHash,
	})
}

// AddChatMembers is used to add members to a chat
func (ctx *Context) AddChatMembers(chatId int64, userIds []int64, forwardLimit int) (bool, error) {
	peerChatStorage := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatStorage.ID == 0 {
		return false, mtp_errors.ErrPeerNotFound
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
		return false, mtp_errors.ErrNotChat
	}
	userPeers := make([]tg.InputUserClass, len(userIds))
	for i, uId := range userIds {
		userPeer := ctx.PeerStorage.GetPeerById(uId)
		if userPeer.ID == 0 {
			return false, mtp_errors.ErrPeerNotFound
		}
		if userPeer.Type != int(storage.TypeUser) {
			return false, mtp_errors.ErrNotUser
		}
		userPeers[i] = &tg.InputUser{
			UserID:     userPeer.ID,
			AccessHash: userPeer.AccessHash,
		}
	}
	return functions.AddChatMembers(ctx, ctx.Raw, chatPeer, userPeers, forwardLimit)
}

// ArchiveChats invokes method folders.editPeerFolders#6847d0ab returning error if any.
// Edit peers in peer folder¹
//
// Links:
//  1. https://core.telegram.org/api/folders#peer-folders
func (ctx *Context) ArchiveChats(chatIds []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIds))
	for i, chatId := range chatIds {
		peer := ctx.PeerStorage.GetPeerById(chatId)
		if peer.ID == 0 {
			return false, mtp_errors.ErrPeerNotFound
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
	return functions.ArchiveChats(ctx, ctx.Raw, chatPeers)
}

// UnarchiveChats invokes method folders.editPeerFolders#6847d0ab returning error if any.
// Edit peers in peer folder¹
//
// Links:
//  1. https://core.telegram.org/api/folders#peer-folders
func (ctx *Context) UnarchiveChats(chatIds []int64) (bool, error) {
	chatPeers := make([]tg.InputPeerClass, len(chatIds))
	for i, chatId := range chatIds {
		peer := ctx.PeerStorage.GetPeerById(chatId)
		if peer.ID == 0 {
			return false, mtp_errors.ErrPeerNotFound
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
func (ctx *Context) CreateChat(title string, userIds []int64) (*tg.Chat, error) {
	userPeers := make([]tg.InputUserClass, len(userIds))
	for i, uId := range userIds {
		userPeer := ctx.PeerStorage.GetPeerById(uId)
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

// DeleteMessages shall be used to delete messages in a chat with chatId and messageIDs.
// Returns error if failed to delete.
func (ctx *Context) DeleteMessages(chatId int64, messageIDs []int) error {
	peer := ctx.PeerStorage.GetPeerById(chatId)
	if peer.ID == 0 {
		return mtp_errors.ErrPeerNotFound
	}
	switch storage.EntityType(peer.Type) {
	case storage.TypeChat:
		_, err := ctx.Raw.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     messageIDs,
		})
		return err
	case storage.TypeChannel:
		_, err := ctx.Raw.ChannelsDeleteMessages(ctx, &tg.ChannelsDeleteMessagesRequest{
			Channel: &tg.InputChannel{
				ChannelID:  peer.ID,
				AccessHash: peer.AccessHash,
			},
			ID: messageIDs,
		})
		return err
	case storage.TypeUser:
		return mtp_errors.ErrNotChat
	default:
		return mtp_errors.ErrPeerNotFound
	}
}

// ForwardMessage shall be used to forward messages in a chat with chatId and messageIDs.
// Returns updatesclass or an error if failed to delete.
//
// Deprecated: use ForwardMessages instead.
func (ctx *Context) ForwardMessage(fromChatId, toChatId int64, request *tg.MessagesForwardMessagesRequest) (tg.UpdatesClass, error) {
	return ctx.ForwardMessages(fromChatId, toChatId, request)
}

// ForwardMessages shall be used to forward messages in a chat with chatId and messageIDs.
// Returns updatesclass or an error if failed to delete.
func (ctx *Context) ForwardMessages(fromChatId, toChatId int64, request *tg.MessagesForwardMessagesRequest) (tg.UpdatesClass, error) {
	fromPeer := ctx.PeerStorage.GetInputPeerById(fromChatId)
	if fromPeer.Zero() {
		return nil, fmt.Errorf("fromChatId: %w", mtp_errors.ErrPeerNotFound)
	}
	toPeer := ctx.PeerStorage.GetInputPeerById(toChatId)
	if toPeer.Zero() {
		return nil, fmt.Errorf("toChatId: %w", mtp_errors.ErrPeerNotFound)
	}
	if request == nil {
		request = &tg.MessagesForwardMessagesRequest{}
	}
	request.RandomID = make([]int64, len(request.ID))
	for i := 0; i < len(request.ID); i++ {
		request.RandomID[i] = ctx.generateRandomID()
	}
	return ctx.Raw.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{
		RandomID: request.RandomID,
		ID:       request.ID,
		FromPeer: fromPeer,
		ToPeer:   toPeer,
	})
}

type EditAdminOpts struct {
	AdminRights tg.ChatAdminRights
	AdminTitle  string
}

// PromoteChatMember is used to promote a user in a chat.
func (ctx *Context) PromoteChatMember(chatId, userId int64, opts *EditAdminOpts) (bool, error) {
	peerChat := ctx.PeerStorage.GetPeerById(chatId)
	if peerChat.ID == 0 {
		return false, fmt.Errorf("chat: %w", mtp_errors.ErrPeerNotFound)
	}
	peerUser := ctx.PeerStorage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return false, fmt.Errorf("user: %w", mtp_errors.ErrPeerNotFound)
	}
	if opts == nil {
		opts = &EditAdminOpts{}
	}
	return functions.PromoteChatMember(ctx, ctx.Raw, peerChat, peerUser, opts.AdminRights, opts.AdminTitle)
}

// DemoteChatMember is used to demote a user in a chat.
func (ctx *Context) DemoteChatMember(chatId, userId int64, opts *EditAdminOpts) (bool, error) {
	peerChat := ctx.PeerStorage.GetPeerById(chatId)
	if peerChat.ID == 0 {
		return false, fmt.Errorf("chat: %w", mtp_errors.ErrPeerNotFound)
	}
	peerUser := ctx.PeerStorage.GetPeerById(userId)
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
			strings.TrimPrefix(
				username,
				"@",
			),
		),
	)
}

func (ctx *Context) extractContactResolvedPeer(p *tg.ContactsResolvedPeer, err error) (types.EffectiveChat, error) {
	if err != nil {
		return &types.EmptyUC{}, err
	}
	go functions.SavePeersFromClassArray(ctx.PeerStorage, p.Chats, p.Users)
	switch p.Peer.(type) {
	case *tg.PeerChannel:
		if p.Chats == nil || len(p.Chats) == 0 {
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
		if p.Users == nil || len(p.Users) == 0 {
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
func (ctx *Context) GetUserProfilePhotos(userId int64, opts *tg.PhotosGetUserPhotosRequest) ([]tg.PhotoClass, error) {
	peerUser := ctx.PeerStorage.GetPeerById(userId)
	if peerUser.ID == 0 {
		return nil, mtp_errors.ErrPeerNotFound
	}
	if opts == nil {
		opts = &tg.PhotosGetUserPhotosRequest{}
	}
	opts.UserID = &tg.InputUser{
		UserID:     userId,
		AccessHash: peerUser.AccessHash,
	}
	p, err := ctx.Raw.PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: opts.UserID,
	})
	if err != nil {
		return nil, err
	}
	return p.GetPhotos(), nil
}

func (ctx *Context) ForwardMediaGroup() error {
	_, err := ctx.Raw.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{})
	return err
}

// ExportSessionString returns session of authorized account in the form of string.
// Note: This session string can be used to log back in with the help of gotgproto.
// Check sessionMaker.SessionType for more information about it.
func (ctx *Context) ExportSessionString() (string, error) {
	return functions.EncodeSessionToString(ctx.PeerStorage.GetSession())
}

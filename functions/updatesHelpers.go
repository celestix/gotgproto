package functions

import (
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func GetNewMessageUpdate(msgData *tg.Message, upds tg.UpdatesClass, p *storage.PeerStorage) *tg.Message {
	u, ok := upds.(*tg.UpdateShortSentMessage)
	if ok {
		msgData.Flags = u.Flags
		msgData.Out = u.Out
		msgData.ID = u.ID
		msgData.Date = u.Date
		msgData.Media = u.Media
		msgData.Entities = u.Entities
		msgData.TTLPeriod = u.TTLPeriod
		return msgData
	}
	for _, update := range GetUpdateClassFromUpdatesClass(upds, p) {
		switch u := update.(type) {
		case *tg.UpdateNewMessage:
			return GetMessageFromMessageClass(u.Message)
		case *tg.UpdateNewChannelMessage:
			return GetMessageFromMessageClass(u.Message)
		case *tg.UpdateNewScheduledMessage:
			return GetMessageFromMessageClass(u.Message)
		}
	}
	return nil
}

func GetEditMessageUpdate(upds tg.UpdatesClass, p *storage.PeerStorage) *tg.Message {
	for _, update := range GetUpdateClassFromUpdatesClass(upds, p) {
		switch u := update.(type) {
		case *tg.UpdateEditMessage:
			return GetMessageFromMessageClass(u.Message)
		case *tg.UpdateEditChannelMessage:
			return GetMessageFromMessageClass(u.Message)
		}
	}
	return nil
}

func GetUpdateClassFromUpdatesClass(updates tg.UpdatesClass, p *storage.PeerStorage) (u []tg.UpdateClass) {
	u, _, _ = getUpdateFromUpdates(updates, p)
	return
}

func getUpdateFromUpdates(updates tg.UpdatesClass, p *storage.PeerStorage) ([]tg.UpdateClass, []tg.ChatClass, []tg.UserClass) {
	switch u := updates.(type) {
	case *tg.Updates:
		go SavePeersFromClassArray(p, u.Chats, u.Users)
		return u.Updates, u.Chats, u.Users
	case *tg.UpdatesCombined:
		go SavePeersFromClassArray(p, u.Chats, u.Users)
		return u.Updates, u.Chats, u.Users
	case *tg.UpdateShort:
		return []tg.UpdateClass{u.Update}, tg.ChatClassArray{}, tg.UserClassArray{}
	default:
		return nil, nil, nil
	}
}

func GetMessageFromMessageClass(m tg.MessageClass) *tg.Message {
	msg, ok := m.(*tg.Message)
	if !ok {
		return nil
	}
	return msg
}

// *************************************************
// *****************INTERNAL-HELPERS****************

func ReturnNewMessageWithError(msgData *tg.Message, upds tg.UpdatesClass, p *storage.PeerStorage, err error) (*tg.Message, error) {
	if err != nil {
		return nil, err
	}
	if msgData == nil {
		msgData = &tg.Message{}
	}
	return GetNewMessageUpdate(msgData, upds, p), nil
}

func ReturnEditMessageWithError(p *storage.PeerStorage, upds tg.UpdatesClass, err error) (*tg.Message, error) {
	if err != nil {
		return nil, err
	}
	return GetEditMessageUpdate(upds, p), nil
}

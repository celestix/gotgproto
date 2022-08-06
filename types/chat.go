package types

import "github.com/gotd/td/tg"

// EffectiveChat interface covers the all three types of chats:
// - tg.User
// - tg.Chat
// - tg.Channel
//
// This interface is implemented by the following structs:
// - User: If the chat is a tg.User then this struct will be returned.
// - Chat: if the chat is a tg.Chat then this struct will be returned.
// - Channel: if the chat is a tg.Channel then this struct will be returned.
// - EmptyUC: if the PeerID doesn't match any of the above cases then EmptyUC struct is returned.
type EffectiveChat interface {
	// Use this method to get chat id.
	GetID() int64
	// Use this method to get access hash of the effective chat.
	GetAccessHash() int64
	// Use this method to check if the effective chat is a channel.
	IsAChannel() bool
	// Use this method to check if the effective chat is a chat.
	IsAChat() bool
	// Use this method to check if the effective chat is a user.
	IsAUser() bool
}

// EmptyUC implements EffectiveChat interface for empty chats.
type EmptyUC struct{}

// Use this method to get chat id.
// Always 0 for EmptyUC
func (*EmptyUC) GetID() int64 {
	return 0
}

// Use this method to get access hash of effective chat.
// Always 0 for EmptyUC
func (*EmptyUC) GetAccessHash() int64 {
	return 0
}

// IsAChannel returns true for a channel.
// Always false for EmptyUC
func (*EmptyUC) IsAChannel() bool {
	return false
}

// IsAChat returns true for a chat.
// Always false for EmptyUC
func (*EmptyUC) IsAChat() bool {
	return false
}

// IsAUser returns true for a user.
// Always false for EmptyUC
func (*EmptyUC) IsAUser() bool {
	return false
}

// User implements EffectiveChat interface for tg.User chats.
type User tg.User

// Use this method to get chat id.
func (u *User) GetID() int64 {
	return u.ID
}

// Use this method to get access hash of the effective chat.
func (u *User) GetAccessHash() int64 {
	return u.AccessHash
}

// IsAChannel returns true for a channel.
func (*User) IsAChannel() bool {
	return false
}

// IsAChat returns true for a chat.
func (*User) IsAChat() bool {
	return false
}

// IsAUser returns true for a user.
func (*User) IsAUser() bool {
	return true
}

func (u *User) Raw() *tg.User {
	us := tg.User(*u)
	return &us
}

// Channel implements EffectiveChat interface for tg.Channel chats.
type Channel tg.Channel

// Use this method to get chat id.
func (u *Channel) GetID() int64 {
	return u.ID
}

// Use this method to get access hash of the effective chat.
func (u *Channel) GetAccessHash() int64 {
	return u.AccessHash
}

// IsAChannel returns true for a channel.
func (*Channel) IsAChannel() bool {
	return true
}

// IsAChat returns true for a chat.
func (*Channel) IsAChat() bool {
	return false
}

// IsAUser returns true for a user.
func (*Channel) IsAUser() bool {
	return false
}

func (u *Channel) Raw() *tg.Channel {
	us := tg.Channel(*u)
	return &us
}

// Chat implements EffectiveChat interface for tg.Chat chats.
type Chat tg.Chat

// Use this method to get chat id.
func (u *Chat) GetID() int64 {
	return u.ID
}

// Use this method to get access hash of the effective chat.
func (*Chat) GetAccessHash() int64 {
	return 0
}

// IsAChannel returns true for a channel.
func (*Chat) IsAChannel() bool {
	return false
}

// IsAChat returns true for a chat.
func (*Chat) IsAChat() bool {
	return true
}

// IsAUser returns true for a user.
func (*Chat) IsAUser() bool {
	return false
}

func (u *Chat) Raw() *tg.Chat {
	us := tg.Chat(*u)
	return &us
}

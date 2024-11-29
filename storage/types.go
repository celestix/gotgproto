package storage

import (
	"errors"

	"github.com/gotd/td/tg"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrDatabaseConnFail = errors.New("failed to connect to database")
	ErrCacheInitFail    = errors.New("failed to initialize cache")
	ErrStorageNotReady  = errors.New("storage is not ready")
)

type EntityType int

func (e EntityType) GetInt() int {
	return int(e)
}

const (
	DefaultUsername   = ""
	DefaultAccessHash = 0
)

const (
	_ EntityType = iota
	TypeUser
	TypeChat
	TypeChannel
)

type Peer struct {
	ID         int64 `gorm:"primary_key"`
	AccessHash int64
	Type       int
	Username   string
}

func getInputPeerFromStoragePeer(peer *Peer) tg.InputPeerClass {
	switch EntityType(peer.Type) {
	case TypeUser:
		return &tg.InputPeerUser{
			UserID:     peer.ID,
			AccessHash: peer.AccessHash,
		}
	case TypeChat:
		return &tg.InputPeerChat{
			ChatID: peer.ID,
		}
	case TypeChannel:
		return &tg.InputPeerChannel{
			ChannelID:  peer.ID,
			AccessHash: peer.AccessHash,
		}
	default:
		return &tg.InputPeerEmpty{}
	}
}

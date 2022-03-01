package storage

import "github.com/gotd/td/tg"

type Peer struct {
	ID         int64 `gorm:"primary_key"`
	AccessHash int64
	Type       int
	Username   string
}

var (
	StoreInMemory = false
	PeerMemoryMap = map[int64]*Peer{}
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

func AddPeer(iD, accessHash int64, peerType EntityType, userName string) {
	peer := &Peer{ID: iD, AccessHash: accessHash, Type: peerType.GetInt(), Username: userName}
	if StoreInMemory {
		PeerMemoryMap[iD] = peer
	} else {
		tx := SESSION.Begin()
		tx.Save(peer)
		tx.Commit()
	}
}

// GetPeerById finds the provided id in the peer storage and return it if found.
func GetPeerById(iD int64) *Peer {
	if StoreInMemory {
		peer := PeerMemoryMap[iD]
		if peer == nil {
			return &Peer{}
		}
		return peer
	} else {
		peer := &Peer{}
		SESSION.Where("id = ?", iD).Find(&peer)
		return peer
	}
}

// GetPeerByUsername finds the provided username in the peer storage and return it if found.
func GetPeerByUsername(username string) *Peer {
	if StoreInMemory {
		for key, peer := range PeerMemoryMap {
			if peer.Username == username {
				return PeerMemoryMap[key]
			}
		}
	} else {
		peer := &Peer{}
		SESSION.Where("username = ?", username).Find(&peer)
		return peer
	}
	return &Peer{}
}

// GetInputPeerById finds the provided id in the peer storage and return its tg.InputPeerClass if found.
func GetInputPeerById(iD int64) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(GetPeerById(iD))
}

// GetInputPeerByUsername finds the provided username in the peer storage and return its tg.InputPeerClass if found.
func GetInputPeerByUsername(userName string) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(GetPeerByUsername(userName))
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

package storage

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"sync"

	"github.com/anonyindian/gotgproto/storage/cache"
	"github.com/gotd/td/tg"
)

type Peer struct {
	ID         int64 `gorm:"primary_key"`
	AccessHash int64
	Type       int
	Username   string
}

var (
	StoreInMemory = false
	peerLock      = sync.RWMutex{}
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
		peerLock.Lock()
		PeerMemoryMap[iD] = peer
		peerLock.Unlock()
	} else {
		go setCachePeers(iD, peer)
		tx := SESSION.Begin()
		tx.Save(peer)
		tx.Commit()
	}
}

// GetPeerById finds the provided id in the peer storage and return it if found.
func GetPeerById(iD int64) *Peer {
	if StoreInMemory {
		peerLock.RLock()
		peer := PeerMemoryMap[iD]
		if peer == nil {
			return &Peer{}
		}
		peerLock.RUnlock()
		return peer
	} else {
		data, err := cache.Cache.Get(strconv.FormatInt(iD, 10))
		if err != nil {
			return cachePeers(iD)
		}
		var peer Peer
		_ = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&peer)
		return &peer
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
		peer := Peer{}
		SESSION.Where("username = ?", username).Find(&peer)
		return &peer
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

func cachePeers(id int64) *Peer {
	var peer = Peer{}
	SESSION.Where("id = ?", id).Find(&peer)
	setCachePeers(id, &peer)
	return &peer
}

func setCachePeers(id int64, peer *Peer) {
	cache.Cache.Set(strconv.FormatInt(id, 10), makeBytes(peer))
}

func makeBytes(v interface{}) []byte {
	buf := bytes.Buffer{}
	_ = gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}

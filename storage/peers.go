package storage

import (
	"sync"

	"github.com/AnimeKaizoku/cacher"
	"github.com/gotd/td/tg"
)

type Peer struct {
	ID         int64 `gorm:"primary_key"`
	AccessHash int64
	Type       int
	Username   string
}

var (
	storeInMemory = false
	peerLock      = sync.RWMutex{}
	peerCache     *cacher.Cacher[int64, *Peer]
	// PeerMemoryMap = map[int64]*Peer{}
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
	if storeInMemory {
		peerCache.Set(iD, peer)
		// PeerMemoryMap[iD] = peer
	} else {
		// go setCachePeers(iD, peer)
		go peerCache.Set(iD, peer)
		tx := SESSION.Begin()
		tx.Save(peer)
		peerLock.Lock()
		defer peerLock.Unlock()
		tx.Commit()
	}
}

// GetPeerById finds the provided id in the peer storage and return it if found.
func GetPeerById(iD int64) *Peer {
	if storeInMemory {
		// peer := PeerMemoryMap[iD]
		peer, ok := peerCache.Get(iD)
		if !ok {
			return &Peer{}
		}
		return peer
	} else {
		peer, ok := peerCache.Get(iD)
		if !ok {
			return cachePeers(iD)
		}
		return peer
		// data := []byte{}
		// data, err := cache.Cache.Get(strconv.FormatInt(iD, 10))
		// if err != nil {
		// 	return cachePeers(iD)
		// }
		// var peer Peer
		// _ = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&peer)
		// return &peer
	}
}

// GetPeerByUsername finds the provided username in the peer storage and return it if found.
func GetPeerByUsername(username string) *Peer {
	if storeInMemory {
		for _, peer := range peerCache.GetAll() {
			if peer.Username == username {
				return peer
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
	// setCachePeers(id, &peer)
	peerCache.Set(id, &peer)
	return &peer
}

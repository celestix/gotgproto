package storage

import (
	"github.com/gotd/td/tg"
)

type Peer struct {
	ID         int64 `gorm:"primary_key"`
	AccessHash int64
	Type       int
	Username   string
}
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

func (p *PeerStorage) AddPeer(iD, accessHash int64, peerType EntityType, userName string) {
	peer := &Peer{ID: iD, AccessHash: accessHash, Type: peerType.GetInt(), Username: userName}
	if p.inMemory {
		p.peerCache.Set(iD, peer)
	} else {
		go p.peerCache.Set(iD, peer)
		tx := p.SqlSession.Begin()
		tx.Save(peer)
		p.peerLock.Lock()
		defer p.peerLock.Unlock()
		tx.Commit()
	}
}

// GetPeerById finds the provided id in the peer storage and return it if found.
func (p *PeerStorage) GetPeerById(iD int64) *Peer {
	if p.inMemory {
		// peer := PeerMemoryMap[iD]
		peer, ok := p.peerCache.Get(iD)
		if !ok {
			return &Peer{}
		}
		return peer
	} else {
		peer, ok := p.peerCache.Get(iD)
		if !ok {
			return p.cachePeers(iD)
		}
		return peer
	}
}

// GetPeerByUsername finds the provided username in the peer storage and return it if found.
func (p *PeerStorage) GetPeerByUsername(username string) *Peer {
	if p.inMemory {
		for _, peer := range p.peerCache.GetAll() {
			if peer.Username == username {
				return peer
			}
		}
	} else {
		peer := Peer{}
		p.SqlSession.Where("username = ?", username).Find(&peer)
		return &peer
	}
	return &Peer{}
}

// GetInputPeerById finds the provided id in the peer storage and return its tg.InputPeerClass if found.
func (p *PeerStorage) GetInputPeerById(iD int64) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(p.GetPeerById(iD))
}

// GetInputPeerByUsername finds the provided username in the peer storage and return its tg.InputPeerClass if found.
func (p *PeerStorage) GetInputPeerByUsername(userName string) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(p.GetPeerByUsername(userName))
}

func (p *PeerStorage) cachePeers(id int64) *Peer {
	var peer = Peer{}
	p.SqlSession.Where("id = ?", id).Find(&peer)
	p.peerCache.Set(id, &peer)
	return &peer
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

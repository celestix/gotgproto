package storage

import (
	"fmt"

	"github.com/gotd/td/tg"
)

func (ps *PeerStorage) AddPeer(iD, accessHash int64, peerType EntityType, userName string) {
	peer := &Peer{
		ID:         iD,
		AccessHash: accessHash,
		Type:       peerType.GetInt(),
		Username:   userName,
	}

	ps.peerCache.Set(iD, peer)

	if !ps.config.Cache.InMemoryOnly {
		go ps.addPeerToDb(peer)
	}
}

func (ps *PeerStorage) addPeerToDb(peer *Peer) {
	ps.peerLock.Lock()
	defer ps.peerLock.Unlock()

	if err := ps.SqlSession.Save(peer).Error; err != nil {
		// TODO: Handle error
	}
}

func (ps *PeerStorage) GetPeerById(iD int64) *Peer {
	peer, ok := ps.peerCache.Get(iD)
	if ps.config.Cache.InMemoryOnly {
		if !ok {
			return &Peer{}
		}
	} else if !ok {
		return ps.cachePeers(iD)
	}

	return peer
}

func (ps *PeerStorage) GetPeerByUsername(username string) *Peer {
	if ps.config.Cache.InMemoryOnly {
		for _, peer := range ps.peerCache.GetAll() {
			if peer.Username == username {
				return peer
			}
		}
	} else {
		var peer Peer
		if err := ps.SqlSession.Where("username = ?", username).First(&peer).Error; err != nil {
			// TODO: Handle error
			return &Peer{}
		}
		return &peer
	}

	return &Peer{}
}

func (ps *PeerStorage) GetInputPeerById(iD int64) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(ps.GetPeerById(iD))
}

func (ps *PeerStorage) GetInputPeerByUsername(userName string) tg.InputPeerClass {
	return getInputPeerFromStoragePeer(ps.GetPeerByUsername(userName))
}

func (ps *PeerStorage) cachePeers(id int64) *Peer {
	var peer Peer

	if err := ps.SqlSession.Where("id = ?", id).First(&peer).Error; err != nil {
		// TODO: Handle error
		return &Peer{}
	}

	ps.peerCache.Set(id, &peer)

	return &peer
}

func (ps *PeerStorage) Close() error {
	if ps.SqlSession != nil {
		db, err := ps.SqlSession.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %w", err)
		}
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	return nil
}

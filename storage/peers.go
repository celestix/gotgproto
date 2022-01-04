package storage

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

const (
	DefaultUsername   = ""
	DefaultAccessHash = 0
	TypeUser          = 1
	TypeChat          = 2
	TypeChannel       = 3
)

func AddPeer(iD, accessHash int64, peerType int, userName string) {
	peer := &Peer{ID: iD, AccessHash: accessHash, Type: peerType, Username: userName}
	if StoreInMemory {
		PeerMemoryMap[iD] = peer
	} else {
		tx := SESSION.Begin()
		tx.Save(peer)
		tx.Commit()
	}
}

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

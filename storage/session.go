package storage

type Session struct {
	Version int `gorm:"primary_key"`
	Data    []byte
}

const LatestVersion = 1

// type Session1 struct {
// 	Version   int `gorm:"primary_key"`
// 	DC        int
// 	Addr      string
// 	AuthKey   []byte
// 	AuthKeyID []byte
// }

func (p *PeerStorage) UpdateSession(session *Session) {
	tx := p.SqlSession.Begin()
	tx.Save(session)
	tx.Commit()
}

// GetSession returns the session saved in storage.
func (p *PeerStorage) GetSession() *Session {
	session := &Session{Version: LatestVersion}
	p.SqlSession.Model(&Session{}).Find(&session)
	return session
}

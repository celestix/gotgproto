package storage

type Session struct {
	DC        int `gorm:"primary_key"`
	Addr      string
	AuthKey   []byte
	AuthKeyID []byte
}

func UpdateSession(session *Session) {
	tx := SESSION.Begin()
	tx.Save(session)
	tx.Commit()
}

func GetSession() *Session {
	session := &Session{}
	SESSION.Model(&Session{}).Find(&session)
	return session
}

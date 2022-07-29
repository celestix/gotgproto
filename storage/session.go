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

func UpdateSession(session *Session) {
	tx := SESSION.Begin()
	tx.Save(session)
	tx.Commit()
}

// GetSession returns the session saved in storage.
func GetSession() *Session {
	session := &Session{Version: LatestVersion}
	SESSION.Model(&Session{}).Find(&session)
	return session
}

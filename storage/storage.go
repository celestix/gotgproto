package storage

import (
	"log"
	"sync"
	"time"

	"github.com/AnimeKaizoku/cacher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PeerStorage struct {
	peerCache  *cacher.Cacher[int64, *Peer]
	peerLock   *sync.RWMutex
	inMemory   bool
	SqlSession *gorm.DB
}

func NewPeerStorage(sessionName string, inMemory bool) *PeerStorage {
	p := PeerStorage{
		inMemory: inMemory,
		peerLock: new(sync.RWMutex),
	}
	var opts *cacher.NewCacherOpts
	if inMemory {
		opts = nil
	} else {
		opts = &cacher.NewCacherOpts{
			TimeToLive:    6 * time.Hour,
			CleanInterval: 24 * time.Hour,
			Revaluate:     true,
		}
		db, err := gorm.Open(sqlite.Open(sessionName), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Panicln(err)
		}
		p.SqlSession = db
		dB, _ := db.DB()
		dB.SetMaxOpenConns(100)
		_ = p.SqlSession.AutoMigrate(&Session{}, &Peer{})
	}
	p.peerCache = cacher.NewCacher[int64, *Peer](opts)
	return &p
}

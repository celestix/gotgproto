package storage

import (
	"log"
	"time"

	"github.com/AnimeKaizoku/cacher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func Load(sessionName string, inMemory bool) {
	loadCache(inMemory)
	if inMemory {
		return
	}
	db, err := gorm.Open(sqlite.Open(sessionName), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panicln(err)
	}
	SESSION = db
	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)

	// Create tables if they don't exist
	_ = SESSION.AutoMigrate(&Session{}, &Peer{})

}

func loadCache(inMemory bool) {
	var opts *cacher.NewCacherOpts
	if inMemory {
		storeInMemory = true
		opts = nil
	} else {
		opts = &cacher.NewCacherOpts{
			TimeToLive:    6 * time.Hour,
			CleanInterval: 24 * time.Hour,
			Revaluate:     true,
		}
	}
	peerCache = cacher.NewCacher[int64, *Peer](opts)
}

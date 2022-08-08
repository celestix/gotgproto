package storage

import (
	"log"

	"github.com/anonyindian/gotgproto/storage/cache"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func Load(sessionName string) {
	err := cache.Load()
	if err != nil {
		log.Panicln(err)
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
	// db.DB().SetMaxOpenConns(100)

	// Create tables if they don't exist
	_ = SESSION.AutoMigrate(&Session{}, &Peer{})

}

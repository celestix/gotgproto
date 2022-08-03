package cache

import (
	"time"

	"github.com/allegro/bigcache"
)

var Cache *bigcache.BigCache

func Load() error {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         1 * time.Hour,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   512,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	var err error
	Cache, err = bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	return nil
}

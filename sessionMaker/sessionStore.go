package sessionMaker

import (
	"context"
	"errors"
	"sync"

	"github.com/amupxm/gotgproto/storage"
	"github.com/gotd/td/session"
)

// SessionStorage implements SessionStorage for file system as file
// stored in Path.
type SessionStorage struct {
	data        []byte
	peerStorage *storage.PeerStorage
	mux         sync.Mutex
}

type jsonData struct {
	Version int
	Data    session.Data
}

// LoadSession loads session from file.
func (f *SessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	if f == nil {
		return nil, errors.New("nil session storage is invalid")
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	return f.data, nil
}

// StoreSession stores session to sqlite storage.
func (f *SessionStorage) StoreSession(_ context.Context, data []byte) error {
	if f == nil {
		return errors.New("nil session storage is invalid")
	}
	f.mux.Lock()
	defer f.mux.Unlock()

	f.peerStorage.UpdateSession(&storage.Session{
		Version: storage.LatestVersion,
		Data:    data,
	})
	return nil
}

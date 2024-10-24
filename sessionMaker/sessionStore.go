package sessionMaker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gotd/td/session"

	"github.com/celestix/gotgproto/storage"
)

// SessionStorage implements SessionStorage for file system as file
// stored in Path.
type SessionStorage struct {
	data        []byte
	peerStorage *storage.PeerStorage
	mux         sync.Mutex
	phone       string
}

type jsonData struct {
	Version int
	Data    session.Data
}

func newSessionStorage(phone string, peerStorage *storage.PeerStorage) *SessionStorage {
	return &SessionStorage{
		phone:       phone,
		peerStorage: peerStorage,
	}
}

// LoadSession loads session from file.
func (f *SessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	if f == nil {
		return nil, errors.New("nil session storage is invalid")
	}

	if d := f.getData(); d != nil {
		return d, nil
	}

	s, err := f.peerStorage.GetSession(f.phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if s == nil {
		return nil, storage.ErrNotFound
	}

	f.setData(s.Data)

	return s.Data, nil
}

// StoreSession stores session to sqlite storage.
func (f *SessionStorage) StoreSession(_ context.Context, data []byte) error {
	if f == nil {
		return errors.New("nil session storage is invalid")
	}

	f.setData(data)

	return f.peerStorage.UpdateSession(&storage.Session{
		Version: storage.LatestVersion,
		Phone:   f.phone,
		Data:    data,
	})
}

func (f *SessionStorage) setData(data []byte) {
	f.mux.Lock()
	f.data = data
	f.mux.Unlock()
}

func (f *SessionStorage) getData() []byte {
	f.mux.Lock()
	d := f.data
	f.mux.Unlock()

	return d
}

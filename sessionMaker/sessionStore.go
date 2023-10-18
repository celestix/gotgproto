package sessionMaker

import (
	"context"
	"sync"

	"github.com/KoNekoD/gotgproto/storage"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
)

// SessionStorage implements SessionStorage for file system as file
// stored in Path.
type SessionStorage struct {
	Session *SessionName
	mux     sync.Mutex
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

	sessionData, err := f.Session.GetData()
	return sessionData, err
	// v := jsonData{
	// 	Version: latestVersion,
	// 	Data:    *sessionData,
	// }
	// data, err := json.Marshal(v)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "marshal")
	// }

	// return data, nil
}

// StoreSession stores session to sqlite storage.
func (f *SessionStorage) StoreSession(_ context.Context, data []byte) error {
	if f == nil {
		return errors.New("nil session storage is invalid")
	}
	f.mux.Lock()
	defer f.mux.Unlock()

	storage.UpdateSession(&storage.Session{
		Version: storage.LatestVersion,
		Data:    data,
	})
	return nil
}

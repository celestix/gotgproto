package sessionMaker

import (
	"context"
	"encoding/json"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"sync"
)

// SessionStorage implements SessionStorage for file system as file
// stored in Path.
type SessionStorage struct {
	Session *SessionName
	mux     sync.Mutex
}

type jsonData struct {
	Version int
	Data    *session.Data
}

const latestVersion = 1

// LoadSession loads session from file.
func (f *SessionStorage) LoadSession(_ context.Context) ([]byte, error) {
	if f == nil {
		return nil, errors.New("nil session storage is invalid")
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	sessionData, err := f.Session.GetData()
	if err != nil {
		return nil, err
	}
	v := jsonData{
		Version: latestVersion,
		Data:    sessionData,
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	return data, nil
}

// StoreSession stores session to sqlite storage.
func (f *SessionStorage) StoreSession(_ context.Context, data []byte) error {
	if f == nil {
		return errors.New("nil session storage is invalid")
	}
	f.mux.Lock()
	defer f.mux.Unlock()
	var v jsonData
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	storage.UpdateSession(&storage.Session{
		DC:        v.Data.DC,
		Addr:      v.Data.Addr,
		AuthKey:   v.Data.AuthKey,
		AuthKeyID: v.Data.AuthKeyID,
	})
	return nil
}

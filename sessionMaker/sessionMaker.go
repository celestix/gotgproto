package sessionMaker

import (
	"context"
	"fmt"

	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

func NewSessionStorage(ctx context.Context, sessionType SessionConstructor, inMemory bool) (*storage.PeerStorage, telegram.SessionStorage, error) {
	name, data, err := sessionType.loadSession()
	if err != nil {
		if err != errSqliteSession {
			return nil, nil, err
		}
		peerStorage := storage.NewPeerStorage(
			fmt.Sprintf("%s.session", name), false)
		return peerStorage, &SessionStorage{
			data:        peerStorage.GetSession().Data,
			peerStorage: peerStorage,
		}, nil
	}
	peerStorage := storage.NewPeerStorage(fmt.Sprintf("%s.session", name), inMemory)
	if inMemory {
		s := session.StorageMemory{}
		err := s.StoreSession(ctx, data)
		if err != nil {
			return nil, nil, err
		}
		return peerStorage, &s, nil
	}
	return peerStorage, &SessionStorage{
		data:        data,
		peerStorage: peerStorage,
	}, nil
}

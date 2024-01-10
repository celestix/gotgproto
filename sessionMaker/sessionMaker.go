package sessionMaker

import (
	"context"
	"fmt"

	"github.com/KoNekoD/gotgproto/storage"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

func NewSessionStorage(
	ctx context.Context,
	sessionType SessionConstructor,
	inMemory bool,
) (*storage.PeerStorage, telegram.SessionStorage, error) {
	name, data, err := sessionType.loadSession()
	if err != nil {
		return nil, nil, err
	}

	if inMemory {
		s := session.StorageMemory{}
		err := s.StoreSession(ctx, data)
		if err != nil {
			return nil, nil, err
		}

		peerStorage := storage.NewPeerStorage("", inMemory)
		return peerStorage, &s, nil
	}

	peerStorage := storage.NewPeerStorage(
		fmt.Sprintf("%s.session", name),
		inMemory,
	)
	s := SessionStorage{data: data, peerStorage: peerStorage}

	return peerStorage, &s, nil
}

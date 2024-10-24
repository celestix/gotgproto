package sessionMaker

import (
	"context"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"

	"github.com/celestix/gotgproto/storage"
)

// NewSessionStorage creates a new session storage based on the provided configuration
func NewSessionStorage(
	ctx context.Context,
	sessionType SessionConstructor,
	phone string,
	cfg *storage.StorageConfig,
) (*storage.PeerStorage, telegram.SessionStorage, error) {
	name, data, err := sessionType.loadSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load session: %w", err)
	}

	var peerStorage *storage.PeerStorage

	// Handle custom dialector case
	if sessDialect, ok := name.(*sessionNameDialector); ok {
		peerStorage, err = storage.NewPeerStorage(
			ctx,
			cfg,
			sessDialect.dialector,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create peer storage with custom dialector: %w", err)
		}

		return peerStorage, newSessionStorage(phone, peerStorage), nil
	}

	// Handle string-based session name
	sessionName := name.(sessionNameString)
	if sessionName == "" {
		sessionName = "gotgproto"
	}

	// Create SQLite-based storage
	dialector := sqlite.Open(fmt.Sprintf("%s.session", sessionName))

	if cfg.Cache.InMemoryOnly {
		// In-memory storage configuration
		peerStorage, err = storage.NewPeerStorage(
			ctx,
			cfg,
			dialector,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create in-memory peer storage: %w", err)
		}

		memStorage := &session.StorageMemory{}
		if err := memStorage.StoreSession(ctx, data); err != nil {
			return nil, nil, fmt.Errorf("failed to store session in memory: %w", err)
		}

		return peerStorage, memStorage, nil
	}

	// Persistent storage configuration
	peerStorage, err = storage.NewPeerStorage(
		ctx,
		cfg,
		dialector,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create persistent peer storage: %w", err)
	}

	return peerStorage, newSessionStorage(phone, peerStorage), nil
}

package sessionMaker

import (
	"fmt"
	"github.com/anonyindian/gotgproto/functions"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/session"
)

type SessionName struct {
	name        string
	sessionType SessionType
}

type SessionType int

const (
	Session SessionType = iota
	StringSession
	TelethonSession
)

func NewSession(sessionName string, sessionType SessionType) *SessionName {
	return &SessionName{
		name:        sessionName,
		sessionType: sessionType,
	}
}

func (s *SessionName) GetName() string {
	return s.name
}

func (s *SessionName) GetData() (*session.Data, error) {
	switch s.sessionType {
	case TelethonSession:
		storage.Load("telethon.session")
		return session.TelethonSession(s.name)
	case StringSession:
		storage.Load("gotgproto.session")
		return functions.DecodeStringToSession(s.name)
	default:
		storage.Load(fmt.Sprintf("%s.session", s.name))
		sFD := storage.GetSession()
		if sFD.DC == 0 {
			return nil, session.ErrNotFound
		}
		return &session.Data{
			DC:        sFD.DC,
			Addr:      sFD.Addr,
			AuthKey:   sFD.AuthKey,
			AuthKeyID: sFD.AuthKeyID,
		}, nil
	}
}

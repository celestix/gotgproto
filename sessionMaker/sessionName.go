package sessionMaker

import (
	"fmt"
	"github.com/anonyindian/gotgproto/functions"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/session"
)

// SessionName object consists of name and SessionType.
type SessionName struct {
	name        string
	sessionType SessionType
}

// SessionType is the type of session you want to log in through.
// It consists of three types: Session, StringSession, TelethonSession.
type SessionType int

const (
	// Session should be used for authorizing into session with default settings.
	Session SessionType = iota
	// StringSession is used as SessionType when you want to log in through the string session made by gotgproto.
	StringSession
	// TelethonSession is used as SessionType when you want to log in through the string session made by telethon - a Python MTProto library.
	TelethonSession
)

// NewSession creates a new session with provided name string and SessionType.
func NewSession(sessionName string, sessionType SessionType) *SessionName {
	return &SessionName{
		name:        sessionName,
		sessionType: sessionType,
	}
}

// GetName is used for retrieving name of the session.
func (s *SessionName) GetName() string {
	return s.name
}

// GetData is used for retrieving session data through provided SessionName type.
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

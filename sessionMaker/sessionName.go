package sessionMaker

import (
	"encoding/json"
	"fmt"

	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/session"
)

// SessionName object consists of name and SessionType.
type SessionName struct {
	name        string
	sessionType SessionType
	data        []byte
	err         error
	PeerStorage *storage.PeerStorage
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
	// PyrogramSession is used as SessionType when you want to log in through the string session made by pyrogram - a Python MTProto library.
	PyrogramSession
)

const (
	InMemorySessionName = ":memory:"
)

// NewSession creates a new session with provided name string and SessionType.
func NewSession(sessionName string, sessionType SessionType) *SessionName {
	s := SessionName{
		name:        sessionName,
		sessionType: sessionType,
	}
	s.data, s.err = s.load()
	return &s
}

// NewInMemorySession Used for create Session with in memory type
func NewInMemorySession(sessionName string, sessionType SessionType) *SessionName {
	s := SessionName{
		name:        InMemorySessionName,
		sessionType: sessionType,
	}
	s.data, s.err = s.loadInMemory(sessionName)
	return &s
}

func (s *SessionName) load() ([]byte, error) {
	switch s.sessionType {
	case PyrogramSession:
		s.PeerStorage = storage.NewPeerStorage("pyrogram.session", false)
		return loadByPyrogramSession(s.name)
	case TelethonSession:
		s.PeerStorage = storage.NewPeerStorage("telethon.session", false)
		return loadByTelethonSession(s.name)
	case StringSession:
		s.PeerStorage = storage.NewPeerStorage("gotgproto.session", false)
		return loadByStringSession(s.name)
	default:
		return s.loadByDefault(s.name)
	}
}

func (s *SessionName) loadInMemory(sessionValue string) ([]byte, error) {
	// ensure that peer caching works properly
	s.PeerStorage = storage.NewPeerStorage("", true)
	switch s.sessionType {
	case PyrogramSession:
		return loadByPyrogramSession(sessionValue)
	case TelethonSession:
		return loadByTelethonSession(sessionValue)
	case StringSession:
		return loadByStringSession(sessionValue)
	case Session:
		panic("can not use sqlite storage with memory storage constructor")
	}
	return nil, nil
}

func loadByPyrogramSession(value string) ([]byte, error) {
	sd, err := DecodePyrogramSession(value)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return data, err
}

func loadByTelethonSession(value string) ([]byte, error) {
	sd, err := session.TelethonSession(value)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return data, err
}

func loadByStringSession(value string) ([]byte, error) {
	sd, err := functions.DecodeStringToSession(value)
	if err != nil {
		return nil, err
	}

	// data, err := json.Marshal(jsonData{
	// 	Version: latestVersion,
	// 	Data:    *sd,
	// })
	return sd.Data, err
}

func (s *SessionName) loadByDefault(value string) ([]byte, error) {
	if value == "" {
		value = "new"
	}
	s.PeerStorage = storage.NewPeerStorage(fmt.Sprintf("%s.session", value), false)
	sFD := s.PeerStorage.GetSession()
	return sFD.Data, nil
}

// GetName is used for retrieving name of the session.
func (s *SessionName) GetName() string {
	return s.name
}

// GetData is used for retrieving session data through provided SessionName type.
func (s *SessionName) GetData() ([]byte, error) {
	return s.data, s.err
}

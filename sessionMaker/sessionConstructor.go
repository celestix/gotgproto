package sessionMaker

import (
	"encoding/json"
	"errors"

	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/session"
)

type SessionConstructor interface {
	loadSession() (string, []byte, error)
}

type SimpleSessionConstructor int8

func SimpleSession() *SimpleSessionConstructor {
	s := SimpleSessionConstructor(0)
	return &s
}

func (s *SimpleSessionConstructor) loadSession() (string, []byte, error) {
	return "gotgproto_simple", nil, nil
}

type SqliteSessionConstructor struct {
	name string
}

func SqliteSession(name string) *SqliteSessionConstructor {
	return &SqliteSessionConstructor{name: name}
}

var errSqliteSession error = errors.New("sqlite session")

func (s *SqliteSessionConstructor) loadSession() (string, []byte, error) {
	if s.name == "" {
		s.name = "new"
	}
	return s.name, nil, errSqliteSession
}

type PyrogramSessionConstructor struct {
	name, value string
}

func PyrogramSession(value string) *PyrogramSessionConstructor {
	return &PyrogramSessionConstructor{value: value}
}

func (s *PyrogramSessionConstructor) Name(name string) *PyrogramSessionConstructor {
	s.name = name
	return s
}

func (s *PyrogramSessionConstructor) loadSession() (string, []byte, error) {
	sd, err := DecodePyrogramSession(s.value)
	if err != nil {
		return s.name, nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return s.name, data, err
}

type TelethonSessionConstructor struct {
	name, value string
}

func TelethonSession(value string) *TelethonSessionConstructor {
	return &TelethonSessionConstructor{value: value}
}

func (s *TelethonSessionConstructor) Name(name string) *TelethonSessionConstructor {
	s.name = name
	return s
}

func (s *TelethonSessionConstructor) loadSession() (string, []byte, error) {
	sd, err := session.TelethonSession(s.value)
	if err != nil {
		return s.name, nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return s.name, data, err
}

type StringSessionConstructor struct {
	name, value string
}

func StringSession(value string) *StringSessionConstructor {
	return &StringSessionConstructor{value: value}
}

func (s *StringSessionConstructor) Name(name string) *StringSessionConstructor {
	s.name = name
	return s
}

func (s *StringSessionConstructor) loadSession() (string, []byte, error) {
	sd, err := functions.DecodeStringToSession(s.value)
	if err != nil {
		return s.name, nil, err
	}
	return s.name, sd.Data, err
}

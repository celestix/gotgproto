package sessionMaker

import (
	"context"
	"encoding/json"

	"github.com/amupxm/gotgproto/functions"
	"github.com/amupxm/gotgproto/storage"
	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"gorm.io/gorm"
)

type sessionName interface {
	getType() string
}

type sessionNameString string

func (sessionNameString) getType() string { return "str" }

type sessionNameDialector struct {
	dialector gorm.Dialector
}

func (sessionNameDialector) getType() string { return "dialector" }

type SessionConstructor interface {
	loadSession() (sessionName, []byte, error)
}

type SimpleSessionConstructor int8

func SimpleSession() *SimpleSessionConstructor {
	s := SimpleSessionConstructor(0)
	return &s
}

func (*SimpleSessionConstructor) loadSession() (sessionName, []byte, error) {
	return sessionNameString("gotgproto_simple"), nil, nil
}

type SqlSessionConstructor struct {
	dialector gorm.Dialector
}

func SqlSession(dialector gorm.Dialector) *SqlSessionConstructor {
	return &SqlSessionConstructor{dialector: dialector}
}

func (s *SqlSessionConstructor) loadSession() (sessionName, []byte, error) {
	return &sessionNameDialector{s.dialector}, nil, nil
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

func (s *PyrogramSessionConstructor) loadSession() (sessionName, []byte, error) {
	sd, err := DecodePyrogramSession(s.value)
	if err != nil {
		return sessionNameString(s.name), nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return sessionNameString(s.name), data, err
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

func (s *TelethonSessionConstructor) loadSession() (sessionName, []byte, error) {
	sd, err := session.TelethonSession(s.value)
	if err != nil {
		return sessionNameString(s.name), nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return sessionNameString(s.name), data, err
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

func (s *StringSessionConstructor) loadSession() (sessionName, []byte, error) {
	sd, err := functions.DecodeStringToSession(s.value)
	if err != nil {
		return sessionNameString(s.name), nil, err
	}
	return sessionNameString(s.name), sd.Data, err
}

type TdataSessionConstructor struct {
	Account tdesktop.Account
	name    string
}

func TdataSession(account tdesktop.Account) *TdataSessionConstructor {
	return &TdataSessionConstructor{Account: account}
}

func (s *TdataSessionConstructor) Name(name string) *TdataSessionConstructor {
	s.name = name
	return s
}

func (s *TdataSessionConstructor) loadSession() (sessionName, []byte, error) {
	sd, err := session.TDesktopSession(s.Account)
	if err != nil {
		return sessionNameString(s.name), nil, err
	}
	ctx := context.Background()
	var (
		gotdstorage = new(session.StorageMemory)
		loader      = session.Loader{Storage: gotdstorage}
	)
	// Save decoded Telegram Desktop session as gotd session.
	if err := loader.Save(ctx, sd); err != nil {
		return sessionNameString(s.name), nil, err
	}
	data, err := json.Marshal(jsonData{
		Version: storage.LatestVersion,
		Data:    *sd,
	})
	return sessionNameString(s.name), data, err
}

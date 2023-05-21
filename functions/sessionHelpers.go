package functions

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/anonyindian/gotgproto/storage"
)

// EncodeSessionToString encodes the provided session to a string in base64 using json bytes.
func EncodeSessionToString(session *storage.Session) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(session)
	if err != nil {
		return "", err
	}
	_ = encoder.Close()
	return buf.String(), nil
}

// DecodeStringToSession decodes the provided base64 encoded session string to session.Data.
func DecodeStringToSession(sessionString string) (*storage.Session, error) {
	// var sessionData session.Data
	var sessionData storage.Session
	return &sessionData, json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(sessionString))).Decode(&sessionData)
}

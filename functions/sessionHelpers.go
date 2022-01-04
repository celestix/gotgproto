package functions

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/session"
	"strings"
)

func EncodeSessionToString(session *storage.Session) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(session)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func DecodeStringToSession(sessionString string) (*session.Data, error) {
	var sessionData session.Data
	return &sessionData, json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(sessionString))).Decode(&sessionData)
}

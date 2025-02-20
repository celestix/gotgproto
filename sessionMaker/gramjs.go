package sessionMaker

import (
	"encoding/base64"
	"encoding/binary"
	"net"
	"strconv"
	"bytes"
	"github.com/go-faster/errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/crypto"
)

func DecodeGramjsSession(hx string) (*session.Data, error) {
	return decodeGramjsSession(hx)
}

func decodeGramjsSession(sessionStr string) (*session.Data, error) {
	data := struct {
		Version       string
		DCID          uint8
		ServerAddress string // 服务器地址字段应为字符串类型
		Port          int16
		Key           []byte
		AuthKey       string
		KeyId         string
	}{}

	if len(sessionStr) == 0 || sessionStr[0] != '1' { // 假设 CURRENT_VERSION 是 "1"
		return nil, errors.New("invalid session string")
	}
	strsession := sessionStr[1:]
	decodedBytes, err := base64.StdEncoding.DecodeString(strsession)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(decodedBytes)

	data.Version = "1" // 固定版本号

	err = binary.Read(reader, binary.BigEndian, &data.DCID)
	if err != nil {
		return nil, err
	}

	addressLength := make([]byte, 2)
	_, err = reader.Read(addressLength)
	if err != nil {
		return nil, err
	}

	addressLen := int(binary.BigEndian.Uint16(addressLength))
	addressBuffer := make([]byte, addressLen)
	_, err = reader.Read(addressBuffer)
	if err != nil {
		return nil, err
	}
	data.ServerAddress = string(bytes.TrimRight(addressBuffer, "\x00")) // 将字节切片转换为字符串并移除可能存在的空字符

	portBuffer := make([]byte, 2)
	_, err = reader.Read(portBuffer)
	if err != nil {
		return nil, err
	}
	data.Port = int16(binary.BigEndian.Uint16(portBuffer))

	keyLength := len(decodedBytes) - (5 + addressLen) // DCID(1 byte), Address Length(2 bytes), Address, Port(2 bytes)
	data.Key = make([]byte, keyLength)
	_, err = reader.Read(data.Key)
	if err != nil {
		return nil, err
	}

	data.AuthKey = base64.StdEncoding.EncodeToString(data.Key)
	keyid := crypto.Key(data.Key).WithID().ID
	data.KeyId = base64.StdEncoding.EncodeToString(keyid[:])
	return &session.Data{
		DC:        int(data.DCID),
		Addr:      net.JoinHostPort(data.ServerAddress, strconv.Itoa(int(data.Port))),
		AuthKey:   data.Key,
		AuthKeyID: keyid[:],
	}, nil
}
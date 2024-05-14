package gotgproto

const (
	clientTypeVPhone int = iota
	clientTypeVBot
)

type clientType interface {
	getType() int
	getValue() string
}

type clientTypePhone string

func (v *clientTypePhone) getType() int {
	return clientTypeVPhone
}

func (v clientTypePhone) getValue() string {
	return string(v)
}

func ClientTypePhone(phoneNumber string) clientType {
	v := clientTypePhone(phoneNumber)
	return &v
}

type clientTypeBot string

func (v *clientTypeBot) getType() int {
	return clientTypeVBot
}

func (v clientTypeBot) getValue() string {
	return string(v)
}

func ClientTypeBot(botToken string) clientType {
	v := clientTypeBot(botToken)
	return &v
}

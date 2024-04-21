package web

import (
	"fmt"

	"github.com/celestix/gotgproto"
)

var authStatus gotgproto.AuthStatus

type webAuth struct{}

var (
	phoneChan  = make(chan string)
	codeChan   = make(chan string)
	passwdChan = make(chan string)
)

func GetWebAuth() gotgproto.AuthConversator {
	return &webAuth{}
}

func (w *webAuth) AskPhoneNumber() (string, error) {
	if authStatus.Event == gotgproto.AuthStatusPhoneRetrial {
		fmt.Println("The phone number you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for phone...")
	code := <-phoneChan
	return code, nil
}

func (w *webAuth) AskCode() (string, error) {
	if authStatus.Event == gotgproto.AuthStatusPhoneCodeRetrial {
		fmt.Println("The OTP you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for code...")
	code := <-codeChan
	return code, nil
}

func (w *webAuth) AskPassword() (string, error) {
	if authStatus.Event == gotgproto.AuthStatusPasswordRetrial {
		fmt.Println("The 2FA password you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for 2fa password...")
	code := <-passwdChan
	return code, nil
}

func (w *webAuth) AuthStatus(authStatusIp gotgproto.AuthStatus) {
	authStatus = authStatusIp
}

func ReceivePhone(phone string) {
	phoneChan <- phone
}

func ReceiveCode(code string) {
	codeChan <- code
}

func ReceivePasswd(passwd string) {
	passwdChan <- passwd
}

package web

import (
	"fmt"

	"github.com/celestix/gotgproto"
)

type webAuth struct {
	phoneChan  chan string
	codeChan   chan string
	passwdChan chan string
	authStatus gotgproto.AuthStatus
}

func GetWebAuth() *webAuth {
	return &webAuth{}
}

func (w *webAuth) AskPhoneNumber() (string, error) {
	if w.authStatus.Event == gotgproto.AuthStatusPhoneRetrial {
		fmt.Println("The phone number you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", w.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for phone...")
	code := <-w.phoneChan
	return code, nil
}

func (w *webAuth) AskCode() (string, error) {
	if w.authStatus.Event == gotgproto.AuthStatusPhoneCodeRetrial {
		fmt.Println("The OTP you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", w.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for code...")
	code := <-w.codeChan
	return code, nil
}

func (w *webAuth) AskPassword() (string, error) {
	if w.authStatus.Event == gotgproto.AuthStatusPasswordRetrial {
		fmt.Println("The 2FA password you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", w.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Println("waiting for 2fa password...")
	code := <-w.passwdChan
	return code, nil
}

func (w *webAuth) AuthStatus(authStatusIp gotgproto.AuthStatus) {
	w.authStatus = authStatusIp
}

func (w *webAuth) ReceivePhone(phone string) {
	w.phoneChan <- phone
}

func (w *webAuth) ReceiveCode(code string) {
	w.codeChan <- code
}

func (w *webAuth) ReceivePasswd(passwd string) {
	w.passwdChan <- passwd
}

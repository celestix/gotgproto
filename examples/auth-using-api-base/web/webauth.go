package web

import (
	"fmt"
	"github.com/celestix/gotgproto"
)

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
	fmt.Println("waiting for phone...")
	code := <-phoneChan
	return code, nil
}

func (w *webAuth) AskCode() (string, error) {
	fmt.Println("waiting for code...")
	code := <-codeChan
	return code, nil
}

func (w *webAuth) AskPassword() (string, error) {
	fmt.Println("waiting for 2fa password...")
	code := <-passwdChan
	return code, nil
}

func (w *webAuth) RetryPassword(attemptsLeft int) (string, error) {
	fmt.Println("The 2FA Code you just entered seems to be incorrect,")
	fmt.Println("Attempts Left:", attemptsLeft)
	fmt.Println("Please try again.... ")
	return w.AskCode()
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

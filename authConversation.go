package gotgproto

import (
	"bufio"
	"fmt"
	"os"
)

type (
	AuthStatusEvent string
	AuthStatus      struct {
		Event        AuthStatusEvent
		AttemptsLeft int
	}
)

func SendAuthStatus(conversator AuthConversator, event AuthStatusEvent) {
	conversator.AuthStatus(AuthStatus{
		Event: event,
	})
}

func SendAuthStatusWithRetrials(conversator AuthConversator, event AuthStatusEvent, attemptsLeft int) {
	conversator.AuthStatus(AuthStatus{
		Event:        event,
		AttemptsLeft: attemptsLeft,
	})
}

var (
	AuthStatusPhoneAsked        = AuthStatusEvent("phone number asked")
	AuthStatusPhoneRetrial      = AuthStatusEvent("phone number validation retrial")
	AuthStatusPhoneFailed       = AuthStatusEvent("phone number validation failed")
	AuthStatusPhoneCodeAsked    = AuthStatusEvent("phone otp asked")
	AuthStatusPhoneCodeVerified = AuthStatusEvent("phone code verified")
	AuthStatusPhoneCodeRetrial  = AuthStatusEvent("phone code verification retrial")
	AuthStatusPhoneCodeFailed   = AuthStatusEvent("phone code verification failed")
	AuthStatusPasswordAsked     = AuthStatusEvent("2fa password asked")
	AuthStatusPasswordRetrial   = AuthStatusEvent("2fa password verification retrial")
	AuthStatusPasswordFailed    = AuthStatusEvent("2fa password verification failed")
	AuthStatusSuccess           = AuthStatusEvent("authentification success")
)

// AuthConversator is an interface for asking user for auth information.
type AuthConversator interface {
	// AskPhoneNumber is called to ask user for phone number.
	// phone number to login should be returned.
	AskPhoneNumber() (string, error)
	// AskCode is called to ask user for OTP.
	// OTP should be returned.
	AskCode() (string, error)
	// AskPassword is called to ask user for 2FA password.
	// 2FA password should be returned.
	AskPassword() (string, error)
	// SendAuthStatus is called to inform the user about
	// the status of the auth process.
	// attemptsLeft is the number of attempts left for the user
	// to enter the input correctly for the current auth status.
	AuthStatus(authStatus AuthStatus)
}

func BasicConversator() AuthConversator {
	return &basicConservator{}
}

type basicConservator struct {
	authStatus AuthStatus
}

func (b *basicConservator) AskPhoneNumber() (string, error) {
	if b.authStatus.Event == AuthStatusPhoneRetrial {
		fmt.Println("The phone number you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", b.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Print("Enter Phone Number: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) AskPassword() (string, error) {
	if b.authStatus.Event == AuthStatusPasswordRetrial {
		fmt.Println("The 2FA password you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", b.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Print("Enter 2FA password: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) AskCode() (string, error) {
	if b.authStatus.Event == AuthStatusPhoneCodeRetrial {
		fmt.Println("The OTP you just entered seems to be incorrect,")
		fmt.Println("Attempts Left:", b.authStatus.AttemptsLeft)
		fmt.Println("Please try again....")
	}
	fmt.Print("Enter Code: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) AuthStatus(authStatus AuthStatus) {
	b.authStatus = authStatus
}

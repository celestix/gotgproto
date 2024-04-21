package gotgproto

import (
	"bufio"
	"fmt"
	"os"
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
	// RetryPassword is called when the 2FA password is incorrect
	// attemptsLeft is the number of attempts left.
	// 2FA password should be returned.
	RetryPassword(attemptsLeft int) (string, error)
}

func BasicConversator() AuthConversator {
	return &basicConservator{}
}

type basicConservator struct{}

func (b *basicConservator) AskPhoneNumber() (string, error) {
	fmt.Print("Enter Phone Number: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) AskPassword() (string, error) {
	fmt.Print("Enter 2FA password: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) AskCode() (string, error) {
	fmt.Print("Enter Code: ")
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (b *basicConservator) RetryPassword(trialsLeft int) (string, error) {
	fmt.Println("The 2FA Code you just entered seems to be incorrect,")
	fmt.Println("Attempts Left:", trialsLeft)
	fmt.Println("Please try again....")
	return b.AskPassword()
}

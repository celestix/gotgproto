package gotgproto

import (
	"bufio"
	"fmt"
	"os"
)

// AuthConversator is an interface for asking user for auth information.
type AuthConversator interface {
	// AskPhoneNumber asks user for phone number.
	AskPhoneNumber() (string, error)
	// AskPassword asks user for 2FA password.
	AskPassword() (string, error)
	// AskCode asks user for auth code.
	AskCode() (string, error)
	// AskRetryCode asks user for retrying 2FA password.
	// attemptsLeft is the number of attempts left.
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
	fmt.Println("Please try again.... ")
	return b.AskCode()
}

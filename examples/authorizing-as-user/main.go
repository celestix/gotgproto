package main

import (
	"context"
	"fmt"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram"
)

func main() {
	// custom dispatcher handles all the updates
	dp := dispatcher.MakeDispatcher()
	gotgproto.StartClient(&gotgproto.ClientHelper{
		// Get AppID from https://my.telegram.org/apps
		AppID: 1234567,
		// Get ApiHash from https://my.telegram.org/apps
		ApiHash: "API_HASH_HERE",
		// Session of your client
		// sessionName: name of the session / session string in case of TelethonSession or StringSession
		// sessionType: can be any out of Session, TelethonSession, StringSession.
		Session: sessionMaker.NewSession("userbot", sessionMaker.Session),
		// Registered Mobile number of the account to be used as the userbot.
		Phone: "PHONE_NUMBER_HERE",
		// Make sure to specify custom dispatcher here in order to enjoy gotgproto's update handling
		Dispatcher: dp,
		// Add the handlers, post functions in TaskFunc
		TaskFunc: func(ctx context.Context, client *telegram.Client) error {
			go func() {
				for {
					if gotgproto.Sender != nil {
						fmt.Println("Client has been started...")
						break
					}
				}
			}()
			return nil
		},
	})
}

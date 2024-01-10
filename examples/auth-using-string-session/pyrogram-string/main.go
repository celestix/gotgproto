package main

import (
	"fmt"
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
)

func main() {
	// Type of client to login to, can be of 2 types:
	// 1.) Bot  (Fill BotToken in this case)
	// 2.) User (Fill Phone in this case)
	clientType := gotgproto.ClientType{
		Phone: "PHONE_NUMBER_HERE",
	}

	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		123456,
		// Get ApiHash from https://my.telegram.org/apps
		"API_HASH_HERE",
		// ClientType, as we defined above
		clientType,
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.PyrogramSession("enter session string here").
				// Sqlite session name (if you're not using memory session)
				// i.e. InMemory in ClientOpts is set to false
				// It will be saved as my_session.session as per this example.
				Name("my_session"),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	fmt.Printf("client (@%s) has been started...\n", client.Self.Username)

	client.Idle()
}

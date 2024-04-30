package main

import (
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/examples/auth-using-api-base/web"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
)

func main() {
	wa := web.GetWebAuth()
	// start web api
	go web.Start(wa)
	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		123456,
		// Get ApiHash from https://my.telegram.org/apps
		"API_HASH_HERE",
		// ClientType, as we defined above
		gotgproto.ClientTypePhone(""),
		// Optional parameters of client
		&gotgproto.ClientOpts{

			// custom authenticator using web api
			AuthConversator: wa,
			Session:         sessionMaker.SqlSession(sqlite.Open("webbot")),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
	client.Idle()

}

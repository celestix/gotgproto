package main

import (
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/examples/auth-using-api-base/web"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
)

func main() {

	// start web api
	go web.Start()

	clientType := gotgproto.ClientType{
		// put your phone here or just leave it empty
		// if you leave it empty, you will be asked to enter your phone number
		Phone: "",
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

			// custom authenticator using web api
			AuthConversator: web.GetWebAuth(),
			Session:         sessionMaker.SqlSession(sqlite.Open("webbot")),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
	client.Idle()

}

package main

import (
	"fmt"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/go-faster/errors"
	"log"
)

func main() {
	// Type of client to login to, can be of 2 types:
	// 1.) Bot  (Fill BotToken in this case)
	// 2.) User (Fill Phone in this case)
	clientType := gotgproto.ClientType{
		BotToken: "BOT_TOKEN_HERE",
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
			InMemory: true,
			Session:  sessionMaker.SimpleSession(),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	clientDispatcher := client.Dispatcher

	// This Message Handler will download any media passed to bot
	clientDispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Media, download), 1)

	fmt.Printf("client (@%s) has been started...\n", client.Self.Username)

	err = client.Idle()
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
}

func download(ctx *ext.Context, update *ext.Update) error {
	filename, err := functions.GetMediaFileNameWithId(update.EffectiveMessage.Media)
	if err != nil {
		return errors.Wrap(err, "failed to get media file name")
	}

	_, err = ctx.DownloadMedia(
		update.EffectiveMessage.Media,
		ext.DownloadOutputPath(filename),
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to download media")
	}

	msg := fmt.Sprintf(`File "%s" downloaded`, filename)
	_, err = ctx.Reply(update, msg, nil)
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	fmt.Println(msg)

	return nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/amupxm/gotgproto/dispatcher/handlers"
	"github.com/amupxm/gotgproto/dispatcher/handlers/filters"
	"github.com/amupxm/gotgproto/ext"
	"github.com/amupxm/gotgproto/functions"
	"github.com/amupxm/gotgproto/sessionMaker"
	"github.com/celestix/gotgproto"
	"github.com/go-faster/errors"
)

func main() {
	appIdEnv := os.Getenv("TG_APP_ID")
	appId, err := strconv.Atoi(appIdEnv)
	if err != nil {
		log.Fatalln("failed to convert app id to int:", err)
	}

	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		appId,
		// Get ApiHash from https://my.telegram.org/apps
		os.Getenv("TG_API_HASH"),
		// ClientType, as we defined above
		gotgproto.ClientTypePhone("PHONE_NUMBER_HERE"),
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

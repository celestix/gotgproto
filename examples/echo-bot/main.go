package main

import (
	"fmt"
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/tg"
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
			Session: sessionMaker.NewSession("echobot", sessionMaker.Session),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	dispatcher := client.Dispatcher

	// Command Handler for /start
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// Callback Query Handler with prefix filter for recieving specific query
	dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), buttonCallback))
	// This Message Handler will call our echo function on new messages
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, echo), 1)

	fmt.Printf("client (@%s) has been started...\n", client.Self.Username)

	client.Idle()
}

// callback function for /start command
func start(ctx *ext.Context, update *ext.Update) error {
	user := update.EffectiveUser()
	_, _ = ctx.Reply(update, fmt.Sprintf("Hello %s, I am @%s and will repeat all your messages.\nI was made using gotd and gotgproto.", user.FirstName, ctx.Self.Username), &ext.ReplyOpts{
		Markup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonURL{
							Text: "gotd/td",
							URL:  "https://github.com/gotd/td",
						},
						&tg.KeyboardButtonURL{
							Text: "gotgproto",
							URL:  "https://github.com/celestix/gotgproto",
						},
					},
				},
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "Click Here",
							Data: []byte("cb_pressed"),
						},
					},
				},
			},
		},
	})
	// End dispatcher groups so that bot doesn't echo /start command usage
	return dispatcher.EndGroups
}

func buttonCallback(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   true,
		QueryID: query.QueryID,
		Message: "This is an example bot!",
	})
	return nil
}

func echo(ctx *ext.Context, update *ext.Update) error {
	msg := update.EffectiveMessage
	_, err := ctx.Reply(update, msg.Text, nil)
	return err
}

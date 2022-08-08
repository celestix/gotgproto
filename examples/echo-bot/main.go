package main

import (
	"context"
	"fmt"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
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
		Session: sessionMaker.NewSession("echobot", sessionMaker.Session),
		// Get BotToken from @botfather
		BotToken: "BOT_TOKEN_HERE",
		// Make sure to specify custom dispatcher here in order to enjoy gotgproto's update handling
		Dispatcher: dp,
		// Add the handlers, post functions in TaskFunc
		TaskFunc: func(ctx context.Context, client *telegram.Client) error {
			// Command Handler for /start
			dp.AddHandler(handlers.NewCommand("start", start))
			// Callback Query Handler with prefix filter for recieving specific query
			dp.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), buttonCallback))
			// This Message Handler will call our echo function on new messages
			dp.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, echo), 1)
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
							URL:  "https://github.com/anonyindian/gotgproto",
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
	_, err := ctx.Reply(update, msg.Message, nil)
	return err
}

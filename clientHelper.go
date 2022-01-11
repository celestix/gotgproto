package gotgproto

import (
	"context"
	"fmt"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/anonyindian/gotgproto/storage"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	// Self is the global variable for the authorized user.
	Self *tg.User
	// Api is the global variable for the tg.Client which is used to make the raw function calls.
	Api *tg.Client
	// Sender is the global variable for message sending helper.
	Sender *message.Sender
)

const VERSION = "v1.0.0-beta03"

type ClientHelper struct {
	// Unique Telegram Application ID, get it from https://my.telegram.org/apps.
	AppID int
	// Unique Telegram API Hash, get it from https://my.telegram.org/apps.
	ApiHash string
	// Session info of the authenticated user, use sessionMaker.NewSession function to fill this field.
	Session *sessionMaker.SessionName
	// BotToken is the unique API Token for the bot you're trying to authorize, get it from @BotFather.
	BotToken string
	// Mobile number of the authenticating user.
	Phone string
	// Dispatcher handlers the incoming updates and execute mapped handlers. It is recommended to use dispatcher.MakeDispatcher function for this field.
	Dispatcher telegram.UpdateHandler
	// TaskFunc is used to for all your post authorization function calls and setting up handlers, check examples for further help.
	TaskFunc func(ctx context.Context, client *telegram.Client) error
	// A Logger provides fast, leveled, structured logging. All methods are safe for concurrent use.
	Logger *zap.Logger
}

// StartClient is the helper for gotd/td which creates client, runs it, prepares storage etc.
func StartClient(c ClientHelper) {
	var sessionStorage telegram.SessionStorage
	if c.Session.GetName() == ":memory:" {
		sessionStorage = &session.StorageMemory{}
		storage.StoreInMemory = true
	} else {
		sessionStorage = &sessionMaker.SessionStorage{
			Session: c.Session,
		}
	}
	c.Run(func(ctx context.Context, log *zap.Logger) error {
		opts := telegram.Options{
			Logger:         c.Logger,
			UpdateHandler:  c.Dispatcher,
			SessionStorage: sessionStorage,
		}
		return c.CreateClient(ctx, opts, c.TaskFunc, telegram.RunUntilCanceled)
	})
	return
}

func (ch ClientHelper) Run(f func(ctx context.Context, log *zap.Logger) error) context.Context {
	clog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = clog.Sync() }()
	ctx := context.Background()
	if err := f(ctx, clog); err != nil {
		clog.Fatal("Run failed", zap.Error(err))
	}
	return ctx
}

func (ch ClientHelper) CreateClient(ctx context.Context, opts telegram.Options,
	setup func(ctx context.Context, Client *telegram.Client) error,
	cb func(ctx context.Context, Client *telegram.Client) error,
) error {
	client := telegram.NewClient(ch.AppID, ch.ApiHash, opts)

	fmt.Printf(`
GoTGProto %s, Copyright (C) 2022 Anony <github.com/anonyindian>
Licensed under the terms of GNU General Public License v3

`, VERSION)

	if err := setup(ctx, client); err != nil {
		return errors.Wrap(err, "setup")
	}

	return client.Run(ctx, func(ctx context.Context) error {
		if ch.BotToken == "" {
			if err := client.Auth().IfNecessary(ctx, auth.NewFlow(termAuth{phone: ch.Phone}, auth.SendCodeOptions{})); err != nil {
				return err
			}
		} else {
			status, err := client.Auth().Status(ctx)
			if err != nil {
				return errors.Wrap(err, "auth status")
			}
			if !status.Authorized {
				if _, err := client.Auth().Bot(ctx, ch.BotToken); err != nil {
					return errors.Wrap(err, "login")
				}
			}
		}
		Self, _ = client.Self(ctx)
		Api = tg.NewClient(client)
		Sender = message.NewSender(Api)
		if ch.Session.GetName() == "" {
			storage.Load("new.session")
		}
		return cb(ctx, client)
	})
}

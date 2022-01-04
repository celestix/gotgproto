package gotgproto

import (
	"context"
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
	Self   *tg.User
	Api    *tg.Client
	Sender *message.Sender
)

type ClientHelper struct {
	AppID      int
	ApiHash    string
	Session    *sessionMaker.SessionName
	BotToken   string
	Phone      string
	Dispatcher telegram.UpdateHandler
	TaskFunc   func(ctx context.Context, client *telegram.Client) error
	Logger     *zap.Logger
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

	if err := setup(ctx, client); err != nil {
		return errors.Wrap(err, "setup")
	}

	return client.Run(ctx, func(ctx context.Context) error {
		if ch.Phone != "" {
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

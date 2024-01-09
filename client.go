package gotgproto

//go:generate go run ./generator

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/celestix/gotgproto/dispatcher"
	intErrors "github.com/celestix/gotgproto/errors"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const VERSION = "v1.0.0-beta14"

type Client struct {
	// Dispatcher handlers the incoming updates and execute mapped handlers. It is recommended to use dispatcher.MakeDispatcher function for this field.
	Dispatcher dispatcher.Dispatcher
	// PublicKeys of telegram.
	//
	// If not provided, embedded public keys will be used.
	PublicKeys []telegram.PublicKey
	// DC ID to connect.
	//
	// If not provided, 2 will be used by default.
	DC int
	// DCList is initial list of addresses to connect.
	DCList dcs.List
	// Resolver to use.
	Resolver dcs.Resolver
	// Whether to show the copyright line in console or no.
	DisableCopyright bool
	// Logger is instance of zap.Logger. No logs by default.
	Logger *zap.Logger

	// Session info of the authenticated user, use sessionMaker.NewSession function to fill this field.
	sessionStorage session.Storage

	// Self contains details of logged in user in the form of *tg.User.
	Self *tg.User

	// Code for the language used on the device's OS, ISO 639-1 standard.
	SystemLangCode string
	// Code for the language used on the client, ISO 639-1 standard.
	ClientLangCode string

	PeerStorage *storage.PeerStorage

	clientType     ClientType
	ctx            context.Context
	err            error
	autoFetchReply bool
	cancel         context.CancelFunc
	running        bool
	*telegram.Client
	appId   int
	apiHash string
}

// Type of client to login to, can be of 2 types:
// 1.) Bot  (Fill BotToken in this case)
// 2.) User (Fill Phone in this case)
type ClientType struct {
	// BotToken is the unique API Token for the bot you're trying to authorize, get it from @BotFather.
	BotToken string
	// Mobile number of the authenticating user.
	Phone string
}

type ClientOpts struct {
	// Logger is instance of zap.Logger. No logs by default.
	Logger *zap.Logger
	// Whether to store session and peer storage in memory or not
	//
	// Note: Sessions and Peers won't be persistent if this field is set to true.
	InMemory bool
	// PublicKeys of telegram.
	//
	// If not provided, embedded public keys will be used.
	PublicKeys []telegram.PublicKey
	// DC ID to connect.
	//
	// If not provided, 2 will be used by default.
	DC int
	// DCList is initial list of addresses to connect.
	DCList dcs.List
	// Resolver to use.
	Resolver dcs.Resolver
	// Whether to show the copyright line in console or no.
	DisableCopyright bool
	// Session info of the authenticated user, use sessionMaker.NewSession function to fill this field.
	Session sessionMaker.SessionConstructor
	// Setting this field to true will lead to automatically fetch the reply_to_message for a new message update.
	//
	// Set to `false` by default.
	AutoFetchReply bool
	// Setting this field to true will lead to automatically fetch the entire reply_to_message chain for a new message update.
	//
	// Set to `false` by default.
	FetchEntireReplyChain bool
	// Code for the language used on the device's OS, ISO 639-1 standard.
	SystemLangCode string
	// Code for the language used on the client, ISO 639-1 standard.
	ClientLangCode string
	// Custom client device
	Device *telegram.DeviceConfig
	// Panic handles all the panics that occur during handler execution.
	PanicHandler dispatcher.PanicHandler
	// Error handles all the unknown errors which are returned by the handler callback functions.
	ErrorHandler dispatcher.ErrorHandler
	// Custom middlewares
	Middlewares []telegram.Middleware
}

// NewClient creates a new gotgproto client and logs in to telegram.
func NewClient(appId int, apiHash string, cType ClientType, opts *ClientOpts) (*Client, error) {
	if opts == nil {
		opts = &ClientOpts{
			SystemLangCode: "en",
			ClientLangCode: "en",
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	peerStorage, sessionStorage, err := sessionMaker.NewSessionStorage(ctx, opts.Session, opts.InMemory)
	if err != nil {
		cancel()
		return nil, err
	}

	// if opts.InMemory {
	// 	d, _ := opts.Session.GetData()
	// 	s := session.StorageMemory{}
	// 	err := s.StoreSession(ctx, d)
	// 	if err != nil {
	// 		cancel()
	// 		return nil, err
	// 	}
	// 	sessionStorage = &s
	// } else {
	// 	sessionStorage = &sessionMaker.SessionStorage{
	// 		Session: opts.Session,
	// 	}
	// }

	d := dispatcher.NewNativeDispatcher(opts.AutoFetchReply, opts.FetchEntireReplyChain, opts.ErrorHandler, opts.PanicHandler, peerStorage)

	// client := telegram.NewClient(appId, apiHash, telegram.Options{
	//	DCList:         opts.DCList,
	//	UpdateHandler: d,
	//	SessionStorage: sessionStorage,
	//	Logger:         opts.Logger,
	//	Device: telegram.DeviceConfig{
	//		DeviceModel:    "GoTGProto",
	//		SystemVersion:  runtime.GOOS,
	//		AppVersion:     VERSION,
	//		SystemLangCode: opts.SystemLangCode,
	//		LangCode:       opts.ClientLangCode,
	//	},
	//	Middlewares: []telegram.Middleware{
	//		floodwait.NewSimpleWaiter().WithMaxRetries(25),
	//		ratelimit.New(rate.Every(100*time.Millisecond), 5),
	//	},
	// })

	c := Client{
		Resolver:         opts.Resolver,
		PublicKeys:       opts.PublicKeys,
		DC:               opts.DC,
		DCList:           opts.DCList,
		DisableCopyright: opts.DisableCopyright,
		Logger:           opts.Logger,
		SystemLangCode:   opts.SystemLangCode,
		ClientLangCode:   opts.ClientLangCode,
		Dispatcher:       d,
		PeerStorage:      peerStorage,
		sessionStorage:   sessionStorage,
		clientType:       cType,
		ctx:              ctx,
		autoFetchReply:   opts.AutoFetchReply,
		cancel:           cancel,
		appId:            appId,
		apiHash:          apiHash,
	}

	c.printCredit()

	return &c, c.Start(opts)
}

func (c *Client) initTelegramClient(
	device *telegram.DeviceConfig,
	middlewares []telegram.Middleware,
) {
	if device == nil {
		device = &telegram.DeviceConfig{
			DeviceModel:    "GoTGProto",
			SystemVersion:  runtime.GOOS,
			AppVersion:     VERSION,
			SystemLangCode: c.SystemLangCode,
			LangCode:       c.ClientLangCode,
		}
	}
	c.Client = telegram.NewClient(c.appId, c.apiHash, telegram.Options{
		DC:             c.DC,
		Resolver:       c.Resolver,
		DCList:         c.DCList,
		UpdateHandler:  c.Dispatcher,
		SessionStorage: c.sessionStorage,
		Logger:         c.Logger,
		Device:         *device,
		Middlewares:    middlewares,
	})
}

func (c *Client) login() error {
	authClient := c.Auth()

	if c.clientType.BotToken == "" {
		authFlow := auth.NewFlow(termAuth{
			phone:  c.clientType.Phone,
			client: authClient,
		},
			auth.SendCodeOptions{})
		if err := IfAuthNecessary(authClient, c.ctx, Flow(authFlow)); err != nil {
			return err
		}
	} else {
		status, err := authClient.Status(c.ctx)
		if err != nil {
			return errors.Wrap(err, "auth status")
		}
		if !status.Authorized {
			if _, err := c.Auth().Bot(c.ctx, c.clientType.BotToken); err != nil {
				return errors.Wrap(err, "login")
			}
		}
	}
	return nil
}

func (ch *Client) printCredit() {
	if !ch.DisableCopyright {
		fmt.Printf(`
GoTGProto %s, Copyright (C) 2023 Anony <github.com/celestix>
Licensed under the terms of GNU General Public License v3

`, VERSION)
	}
}

func (c *Client) initialize(wg *sync.WaitGroup) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := c.login()
		if err != nil {
			return err
		}
		self, err := c.Client.Self(ctx)
		if err != nil {
			return err
		}

		c.Self = self

		c.Dispatcher.Initialize(ctx, c.Stop, c.Client, self)

		c.PeerStorage.AddPeer(self.ID, self.AccessHash, storage.TypeUser, self.Username)
		// notify channel that client is up
		wg.Done()
		c.running = true
		<-c.ctx.Done()
		return c.ctx.Err()
	}
}

// ExportStringSession EncodeSessionToString encodes the client session to a string in base64.
//
// Note: You must not share this string with anyone, it contains auth details for your logged in account.
func (c *Client) ExportStringSession() (string, error) {
	// InMemorySession case
	loadSession, err := c.sessionStorage.LoadSession(c.ctx)
	if err == nil {
		return string(loadSession), nil
	}

	// todo. what if session is InMemorySession? We got panic
	return functions.EncodeSessionToString(c.PeerStorage.GetSession())
}

// Idle keeps the current goroutined blocked until the client is stopped.
func (c *Client) Idle() error {
	<-c.ctx.Done()
	return c.err
}

// CreateContext creates a new pseudo updates context.
// A context retrieved from this method should be reused.
func (c *Client) CreateContext() *ext.Context {
	return ext.NewContext(
		c.ctx,
		c.API(),
		c.PeerStorage,
		c.Self,
		message.NewSender(c.API()),
		&tg.Entities{
			Users: map[int64]*tg.User{
				c.Self.ID: c.Self,
			},
		},
		c.autoFetchReply,
	)
}

// Stop cancels the context.Context being used for the client
// and stops it.
//
// Notes:
//
// 1.) Client.Idle() will exit if this method is called.
//
// 2.) You can call Client.Start() to start the client again
// if it was stopped using this method.
func (c *Client) Stop() {
	c.cancel()
	c.running = false
}

// Start connects the client to telegram servers and logins.
// It will return error if the client is already running.
func (c *Client) Start(opts *ClientOpts) error {
	if c.running {
		return intErrors.ErrClientAlreadyRunning
	}
	if c.ctx.Err() == context.Canceled {
		c.ctx, c.cancel = context.WithCancel(context.Background())
	}
	c.initTelegramClient(opts.Device, opts.Middlewares)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c *Client) {
		c.err = c.Run(c.ctx, c.initialize(&wg))
		if c.err != nil {
			wg.Done()
		}
	}(c)

	// wait till client starts
	wg.Wait()
	return c.err
}

// RefreshContext casts the new context.Context and telegram session
// to ext.Context (It may be used after doing Stop and Start calls respectively.)
func (c *Client) RefreshContext(ctx *ext.Context) {
	(*ctx).Context = c.ctx
	(*ctx).Raw = c.API()
}

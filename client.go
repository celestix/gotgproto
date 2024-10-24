package gotgproto

//go:generate go run ./generator

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"

	"github.com/celestix/gotgproto/dispatcher"
	intErrors "github.com/celestix/gotgproto/errors"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/celestix/gotgproto/storage"
)

const VERSION = "v1.0.0-beta18"

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
	// MigrationTimeout configures migration timeout.
	MigrationTimeout time.Duration
	// AckBatchSize is limit of MTProto ACK buffer size.
	AckBatchSize int
	// AckInterval is maximum time to buffer MTProto ACK.
	AckInterval time.Duration
	// RetryInterval is duration between send retries.
	RetryInterval time.Duration
	// MaxRetries is limit of send retries.
	MaxRetries int
	// ExchangeTimeout is timeout of every key exchange request.
	ExchangeTimeout time.Duration
	// DialTimeout is timeout of creating connection.
	DialTimeout time.Duration
	// CompressThreshold is a threshold in bytes to determine that message
	// is large enough to be compressed using GZIP.
	// If < 0, compression will be disabled.
	// If == 0, default value will be used.
	CompressThreshold int
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
	// PeerStorage is the storage for all the peers.
	// It is recommended to use storage.NewPeerStorage function for this field.
	PeerStorage *storage.PeerStorage
	// NoAutoAuth is a flag to disable automatic authentication
	// if the current session is invalid.
	NoAutoAuth bool

	*telegram.Client
	authConversator AuthConversator
	clientType      clientType
	ctx             context.Context
	err             error
	autoFetchReply  bool
	cancel          context.CancelFunc
	running         bool
	appId           int
	apiHash         string
	opts            *ClientOpts
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
	// StorageConfig is the configuration for the storage.
	StorageConfig *storage.StorageConfig
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
	// Custom Middlewares
	Middlewares []telegram.Middleware
	// Custom Run() Middleware
	// Can be used for floodWaiter package
	// https://github.com/celestix/gotgproto/blob/beta/examples/middleware/main.go#L41
	RunMiddleware func(
		origRun func(ctx context.Context, f func(ctx context.Context) error) (err error),
		ctx context.Context,
		f func(ctx context.Context) (err error),
	) (err error)
	// A custom context to use for the client.
	// If not provided, context.Background() will be used.
	// Note: This context will be used for the entire lifecycle of the client.
	Context context.Context
	// AuthConversator is the interface for the authenticator.
	// gotgproto.BasicConversator is used by default.
	AuthConversator AuthConversator
	// MigrationTimeout configures migration timeout.
	MigrationTimeout time.Duration
	// AckBatchSize is limit of MTProto ACK buffer size.
	AckBatchSize int
	// AckInterval is maximum time to buffer MTProto ACK.
	AckInterval time.Duration
	// RetryInterval is duration between send retries.
	RetryInterval time.Duration
	// MaxRetries is limit of send retries.
	MaxRetries int
	// ExchangeTimeout is timeout of every key exchange request.
	ExchangeTimeout time.Duration
	// DialTimeout is timeout of creating connection.
	DialTimeout time.Duration
	// CompressThreshold is a threshold in bytes to determine that message
	// is large enough to be compressed using GZIP.
	// If < 0, compression will be disabled.
	// If == 0, default value will be used.
	CompressThreshold int
	// NoAutoAuth is a flag to disable automatic authentication
	// if the current session is invalid.
	NoAutoAuth bool
	// NoAutoStart is a flag to disable automatic start of the client.
	NoAutoStart bool
}

var DefaultOpts = &ClientOpts{
	SystemLangCode: "en",
	ClientLangCode: "en",
}

// NewClient creates a new gotgproto client and logs in to telegram.
func NewClient(appId int, apiHash string, cType clientType, opts *ClientOpts) (*Client, error) {
	if opts == nil {
		opts = DefaultOpts
	}

	if opts.StorageConfig == nil {
		opts.StorageConfig = storage.DefaultStorageConfig
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}
	ctx, cancel := context.WithCancel(opts.Context)

	peerStorage, sessionStorage, err := sessionMaker.NewSessionStorage(ctx, opts.Session, cType.getValue(), opts.StorageConfig)
	if err != nil {
		cancel()
		return nil, err
	}

	if opts.AuthConversator == nil {
		opts.AuthConversator = BasicConversator()
	}

	d := dispatcher.NewNativeDispatcher(opts.AutoFetchReply, opts.FetchEntireReplyChain, opts.ErrorHandler, opts.PanicHandler, peerStorage)

	c := Client{
		Resolver:          opts.Resolver,
		PublicKeys:        opts.PublicKeys,
		DC:                opts.DC,
		DCList:            opts.DCList,
		MigrationTimeout:  opts.MigrationTimeout,
		AckBatchSize:      opts.AckBatchSize,
		AckInterval:       opts.AckInterval,
		RetryInterval:     opts.RetryInterval,
		MaxRetries:        opts.MaxRetries,
		ExchangeTimeout:   opts.ExchangeTimeout,
		DialTimeout:       opts.DialTimeout,
		CompressThreshold: opts.CompressThreshold,
		DisableCopyright:  opts.DisableCopyright,
		Logger:            opts.Logger,
		SystemLangCode:    opts.SystemLangCode,
		ClientLangCode:    opts.ClientLangCode,
		NoAutoAuth:        opts.NoAutoAuth,
		authConversator:   opts.AuthConversator,
		Dispatcher:        d,
		PeerStorage:       peerStorage,
		sessionStorage:    sessionStorage,
		clientType:        cType,
		ctx:               ctx,
		autoFetchReply:    opts.AutoFetchReply,
		cancel:            cancel,
		appId:             appId,
		apiHash:           apiHash,
		opts:              opts,
	}

	c.printCredit()

	if opts.NoAutoStart {
		return &c, nil
	}

	return &c, c.Start()
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
		DCList:            c.DCList,
		Resolver:          c.Resolver,
		DC:                c.DC,
		PublicKeys:        c.PublicKeys,
		MigrationTimeout:  c.MigrationTimeout,
		AckBatchSize:      c.AckBatchSize,
		AckInterval:       c.AckInterval,
		RetryInterval:     c.RetryInterval,
		MaxRetries:        c.MaxRetries,
		ExchangeTimeout:   c.ExchangeTimeout,
		DialTimeout:       c.DialTimeout,
		CompressThreshold: c.CompressThreshold,
		UpdateHandler:     c.Dispatcher,
		SessionStorage:    c.sessionStorage,
		Logger:            c.Logger,
		Device:            *device,
		Middlewares:       middlewares,
	})
}

func (c *Client) Login(conversator ...AuthConversator) error {
	authClient := c.Auth()

	status, err := authClient.Status(c.ctx)
	if err != nil {
		return fmt.Errorf("auth status: %w", err)
	}

	if status.Authorized {
		return nil
	}

	_conversator := c.authConversator
	if len(conversator) > 0 {
		_conversator = conversator[0]
	}

	switch c.clientType.getType() {
	case clientTypeVPhone:
		if c.NoAutoAuth {
			return intErrors.ErrSessionUnauthorized
		}

		phoneNr := c.clientType.getValue()

		if err := newAuthFlow(
			authClient,
			_conversator,
			phoneNr,
			auth.SendCodeOptions{},
		).Execute(c.ctx); err != nil {
			return fmt.Errorf("auth flow: %w", err)
		}
	case clientTypeVBot:
		if !status.Authorized {
			if _, err := c.Auth().Bot(c.ctx, c.clientType.getValue()); err != nil {
				return fmt.Errorf("bot auth: %w", err)
			}
		}
	default:
		return fmt.Errorf("invalid client type, must be either clientTypeVPhone or clientTypeVBot")
	}

	return nil
}

func (ch *Client) printCredit() {
	if !ch.DisableCopyright {
		fmt.Printf(`
GoTGProto %s, Copyright (C) 2024 Anony <github.com/celestix>
Licensed under the terms of GNU General Public License v3

`, VERSION)
	}
}

func (c *Client) initialize(wg *sync.WaitGroup) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := c.Login(); err != nil {
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
	loadedSessionData, err := c.sessionStorage.LoadSession(c.ctx)
	if err == nil {
		loadedSession := &storage.Session{
			Version: storage.LatestVersion,
			Data:    loadedSessionData,
		}
		return functions.EncodeSessionToString(loadedSession)
	}

	session, err := c.PeerStorage.GetSession(c.clientType.getValue())
	if err != nil {
		return "", fmt.Errorf("get session: %w", err)
	}

	return functions.EncodeSessionToString(session)
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
func (c *Client) Start(opts ...*ClientOpts) error {
	if c.running {
		return intErrors.ErrClientAlreadyRunning
	}

	if len(opts) > 0 {
		c.opts = opts[0]
	}

	if c.opts == nil {

	}

	if c.appId == 0 || len(c.apiHash) == 0 {
		return intErrors.ErrClientNotInitialized
	}

	if c.ctx.Err() == context.Canceled {
		c.ctx, c.cancel = context.WithCancel(context.Background())
	}

	c.initTelegramClient(c.opts.Device, c.opts.Middlewares)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(c *Client) {
		if c.opts.RunMiddleware == nil {
			c.err = c.Run(c.ctx, c.initialize(&wg))
		} else {
			c.err = c.opts.RunMiddleware(
				c.Run,
				c.ctx,
				c.initialize(&wg),
			)
		}

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

# <a href="https://github.com/celestix/gotgproto"><img src="./gotgproto.png" width="40px" align="left"></img></a> GoTGProto
GoTGProto is a helper package for gotd library, It aims to make td's raw functions easy-to-use with the help of features like using session strings, custom helper functions, storing peers and extracting chat or user ids through it etc.

We have an outstanding userbot project going on with GoTGProto, you can check it out by [clicking here](https://github.com/GigaUserbot/GIGA). 

You can use this package to create bots and userbots with Telegram MTProto easily in golang, for any futher help you can check out the [documentations](https://pkg.go.dev/github.com/celestix/gotgproto) or reach us through the following:
- Updates Channel: [![Channel](https://img.shields.io/badge/GoTGProto-Channel-dark)](https://telegram.me/gotgproto)
- Support Chat: [![Chat](https://img.shields.io/badge/GoTGProto-Support%20Chat-red)](https://telegram.me/gotgprotochat)

[![Go Reference](https://pkg.go.dev/badge/github.com/celestix/gotgproto.svg)](https://pkg.go.dev/github.com/celestix/gotgproto) [![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](http://perso.crans.org/besson/LICENSE.html)

**Note**: This library is in the beta stage yet and may not be stable for every case.

## Installation
You can download the library with the help of standard `go get` command.

```bash
go get github.com/celestix/gotgproto
```

## Usage
You can find various examples in the [examples' directory](./examples/), one of them i.e. authorizing as a user is as follows:
```go
package main

import (
	"log"
	
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
)

func main() {
	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		123456,
		// Get ApiHash from https://my.telegram.org/apps
		"API_HASH_HERE",
		// ClientType, as we defined above
		gotgproto.ClientTypePhone("PHONE_NUMBER_HERE"),
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("echobot")),
		},
	)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
	client.Idle()
}
```

## Basic Operations
Here are some quick examples on basic operations like sending a message, media etc.

Naming convention:
- `ctx` is a `*ext.Context` object returned as a paramter in all update handlers.
- `update` is a `*ext.Update` object returned as a parameter in all update handlers.
- `chatId` is the chat id of the chat you want to send the message to. (type int64)

**Note**: You do not need to specify the peer field in the request, it is automatically filled by the library.

### Sending a Message
```go
ctx.SendMessage(chatId, &tg.MessagesSendMessageRequest{
		Message: "Hello, World!",
		// Peer: ... (No need of setting peer as we have passed chatId)
})
```

### Uploading media on telegram
If you want to send a local file, you will need to upload it to telegram using an uploader instance as we've done below for `test.jpg`:
```go
f, err := uploader.NewUploader(ctx.Raw).FromPath(ctx, "test.jpg")
if err != nil {
	panic(err)
}
```

### Sending uploaded media to a chat
Let's upload the photo (`test.jpg`) we just uploaded on telegram:
```go
ctx.SendMedia(chatId, &tg.MessagesSendMediaRequest{
		Message: "This is your caption",
		Media: &tg.InputMediaUploadedPhoto{
			File: f,
			MimeType: "video/mp4", // You can specify the mime type of the file here like "image/jpeg",  "audio/mpeg" etc.
			Attributes: []tg.DocumentAttributeClass{&tg.DocumentAttributeFilename{FileName: f.GetName()}}
		},
})
```

_For media types other than photos, use `tg.InputMediaUploadedDocument`._

### Retrieving a photo from a message and sending it
If you want to send a photo from a message, you can do it like this:
```go
m := update.EffectiveMessage
// we recommend you to check if the media is a photo casting it in real life applications.
photo := m.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo)
ctx.SendMedia(chatId, &tg.MessagesSendMediaRequest{
	Media: &tg.InputMediaPhoto{
		// Specifying ID, AccessHash and FileReference of the photo is compulsory.
		ID: &tg.InputPhoto{
			ID:            photo.ID,
			AccessHash:    photo.AccessHash,
			FileReference: photo.FileReference,
		},
	},
})
```

### Sending a file to a chat after retrieving it from a message
```go
m := update.EffectiveMessage
// we recommend you to check if the media is a photo casting it in real life applications.
doc := m.Media.(*tg.MessageMediaDocument).Document.(*tg.Document)
ctx.SendMedia(chatId, &tg.MessagesSendMediaRequest{
	Media: &tg.InputMediaDocument{
		ID: &tg.InputDocument{
			ID:            doc.ID,
			AccessHash:    doc.AccessHash,
			FileReference: doc.FileReference,
		},
	},
})
```

### Working with raw tl functions
Telegram has a big library of functions, Gotgproto doesn't have helper for all of them currently, but you can use the raw functions to call any function you want and also utilize this library's features. Here is an example of calling the `messages.getHistory` function to get chat history:
```go
// peer storage is managed by the library automatically with each session. It stores the chat ids and their access hash which are needed to create input peer queries.
peerStorage = ctx.PeerStorage
// get the peer from the chat id
inputPeer := peerStorage.GetInputPeerById(chatId)
// draw out a raw function call using ctx.Raw api
ctx.Raw.MessagesGetHistory(
	ctx,
	&tg.MessagesGetHistoryRequest{
		// Peer is compulsory
		Peer:  inputPeer,
		Limit: 10,
	},
)
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update the examples as appropriate.

## License
[![GPLv3](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html)
<br>Licensed Under <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License v3</a>

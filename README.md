# GoTGProto
GoTGProto is a helper package for gotd library, It aims to make td's raw functions easy-to-user with the help of features like using session strings, custom helper functions, storing peers and extracting chat or user ids through it etc.

You can use this package to create bots and userbots with Telegram MTProto easily in golang, for any futher help you can check out the [documentations](https://pkg.go.dev/github.com/anonyindian/gotgproto) or reach us through the following:
- Updates Channel: [![Channel](https://img.shields.io/badge/GoTGProto-Channel-dark)](https://telegram.me/gotgproto)
- Support Chat: [![Chat](https://img.shields.io/badge/GoTGProto-Support%20Chat-red)](https://telegram.me/gotgprotochat)

[![Go Reference](https://pkg.go.dev/badge/github.com/anonyindian/gotgproto.svg)](https://pkg.go.dev/github.com/anonyindian/gotgproto) [![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](http://perso.crans.org/besson/LICENSE.html)

**Note**: This library is in the beta stage yet and may not be stable for every case.

## Installation
You can download the library with the help of standard `go get` command.

```bash
go get github.com/anonyindian/gotgproto
```

## Usage
You can find various examples in the [examples' directory](./examples/), one of them i.e. authorizing as a user is as follows:
```go
package main
import (
	"context"
	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram"
)
func main() {
    dp := dispatcher.MakeDispatcher()
    gotgproto.StartClient(gotgproto.ClientHelper{
		AppID: 1234567,
		ApiHash: "API_HASH_HERE",
		Session: sessionMaker.NewSession("session_name", sessionMaker.Session),
		Phone: "PHONE_NUMBER_HERE",
		Dispatcher: dp,
		TaskFunc: func(ctx context.Context, client *telegram.Client) error {
			return nil
		},
	})
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update examples as appropriate.

## License
[![GPLv3](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html)
<br>Licensed Under <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License v3</a>
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
	"github.com/amupxm/gotgproto/sessionMaker"
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

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update the examples as appropriate.

## License
[![GPLv3](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html)
<br>Licensed Under <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License v3</a>

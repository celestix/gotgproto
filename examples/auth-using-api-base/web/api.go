package web

import (
	"fmt"
	"net/http"

	"github.com/celestix/gotgproto"
)

// Start a web server and wait
func Start() {
	http.HandleFunc("/", setInfo)
	http.HandleFunc("/getAuthStatus", getAuthStatus)
	http.ListenAndServe(":9997", nil)
}

func getAuthStatus(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, authStatus.Event)
}

// setInfo handle user info, set phone, code or passwd
func setInfo(w http.ResponseWriter, req *http.Request) {
	action := req.URL.Query().Get("set")

	switch action {

	case "phone":
		fmt.Println("Rec phone")
		num := req.URL.Query().Get("phone")
		phone := "+" + num
		ReceivePhone(phone)
		for authStatus.Event == gotgproto.AuthStatusPhoneAsked ||
			authStatus.Event == gotgproto.AuthStatusPhoneRetrial {
			continue
		}
	case "code":
		fmt.Println("Rec code")
		code := req.URL.Query().Get("code")
		ReceiveCode(code)
		for authStatus.Event == gotgproto.AuthStatusPhoneCodeAsked ||
			authStatus.Event == gotgproto.AuthStatusPhoneCodeRetrial {
			continue
		}
	case "passwd":
		passwd := req.URL.Query().Get("passwd")
		ReceivePasswd(passwd)
		for authStatus.Event == gotgproto.AuthStatusPasswordAsked ||
			authStatus.Event == gotgproto.AuthStatusPasswordRetrial {
			continue
		}
	}
	w.Write([]byte(authStatus.Event))
}

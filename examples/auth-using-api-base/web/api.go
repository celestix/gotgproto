package web

import (
	"fmt"
	"net/http"

	"github.com/celestix/gotgproto"
)

// Start a web server and wait
func Start(wa *webAuth) {
	http.HandleFunc("/", wa.setInfo)
	http.HandleFunc("/getAuthStatus", wa.getAuthStatus)
	http.ListenAndServe(":9997", nil)
}

func (wa *webAuth) getAuthStatus(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, wa.authStatus.Event)
}

// setInfo handle user info, set phone, code or passwd
func (wa *webAuth) setInfo(w http.ResponseWriter, req *http.Request) {
	action := req.URL.Query().Get("set")

	switch action {

	case "phone":
		fmt.Println("Rec phone")
		num := req.URL.Query().Get("phone")
		phone := "+" + num
		wa.ReceivePhone(phone)
		for wa.authStatus.Event == gotgproto.AuthStatusPhoneAsked ||
			wa.authStatus.Event == gotgproto.AuthStatusPhoneRetrial {
			continue
		}
	case "code":
		fmt.Println("Rec code")
		code := req.URL.Query().Get("code")
		wa.ReceiveCode(code)
		for wa.authStatus.Event == gotgproto.AuthStatusPhoneCodeAsked ||
			wa.authStatus.Event == gotgproto.AuthStatusPhoneCodeRetrial {
			continue
		}
	case "passwd":
		passwd := req.URL.Query().Get("passwd")
		wa.ReceivePasswd(passwd)
		for wa.authStatus.Event == gotgproto.AuthStatusPasswordAsked ||
			wa.authStatus.Event == gotgproto.AuthStatusPasswordRetrial {
			continue
		}
	}
	w.Write([]byte(wa.authStatus.Event))
}

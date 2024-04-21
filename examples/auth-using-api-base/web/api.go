package web

import (
	"fmt"
	"net/http"
)

// Start a web server and wait
func Start() {
	http.HandleFunc("/", setInfo)
	http.HandleFunc("/getAuthStatus", getAuthStatus)
	http.ListenAndServe(":9997", nil)
}

func getAuthStatus(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, authStatus)
}

// setInfo handle user info, set phone, code or passwd
func setInfo(w http.ResponseWriter, req *http.Request) {
	action := req.URL.Query().Get("set")

	switch action {

	case "phone":
		num := req.URL.Query().Get("phone")
		phone := "+" + num
		ReceivePhone(phone)
		fmt.Fprintf(w, "phone received: %s", phone)

	case "code":
		code := req.URL.Query().Get("code")
		ReceiveCode(code)
		fmt.Fprintf(w, "code received: %s", code)

	case "passwd":
		passwd := req.URL.Query().Get("passwd")
		ReceivePasswd(passwd)
		fmt.Fprintf(w, "passwd received: %s", passwd)

	}
}

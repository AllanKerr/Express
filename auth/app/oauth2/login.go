package oauth2

import (
	"net/http"
	"html/template"
	"fmt"
	"net/url"
)

func (ctrl *HTTPController) Login(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, nil)
}

func (ctrl *HTTPController) Submit(w http.ResponseWriter, req *http.Request) {

	username := req.FormValue("username")
	password := req.FormValue("password")

	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", "wlYVsgnUibXbqMtK-f0aXgXJQiHdOucU47uGUg48Zx8=")
	form.Add("client_secret", "EGWWqWWLXyRdOp3NowbPDj8YURwWPfWcGtEyD6cVk2s=")
	form.Add("username", username)
	form.Add("password", password)
	req.PostForm = form

	ctrl.Token(w, req)
}

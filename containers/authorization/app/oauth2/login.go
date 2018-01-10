package oauth2

import (
	"net/http"
	"html/template"
	"fmt"
	"net/url"
	"os"
)

// HTTP GET handler for displaying the login page
func (ctrl *HTTPController) Login(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, nil)
}

// HTTP post handler for handling login submission
func (ctrl *HTTPController) SubmitLogin(w http.ResponseWriter, req *http.Request) {

	username := req.FormValue("username")
	password := req.FormValue("password")

	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", os.Getenv("CLIENT_ID"))
	form.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	form.Add("scope", "user offline")
	form.Add("username", username)
	form.Add("password", password)
	req.PostForm = form

	ctrl.Token(w, req)
}

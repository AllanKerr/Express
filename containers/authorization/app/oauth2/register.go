package oauth2

import (
	"fmt"
	"html/template"
	"net/http"
	"github.com/ory/fosite"
)

func (ctrl *HTTPController) Register(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, nil)
}

func (ctrl *HTTPController) SubmitRegistration(w http.ResponseWriter, req *http.Request) {

	username := req.FormValue("username")
	password := req.FormValue("password")
	confirmPassword := req.FormValue("confirm-password")

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := NewUser(username, password)
	err := ctrl.adapter.CreateUser(user)
	if err == fosite.ErrInvalidRequest {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctrl.Submit(w, req)
}

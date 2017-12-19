package oauth2

import (
	"net/http"
	"html/template"
	"fmt"
	"os"
)

func (ctrl *HTTPController) Login(w http.ResponseWriter, req *http.Request) {

	cwd, _ := os.Getwd()
	fmt.Println("CWD: " + cwd)

	t, err := template.ParseFiles("templates/welcome.html") //setp 1
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, "Hello World!") //step 2
}

package http

import (
	"errors"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"rsvp/pkg/rsvp"
)

var (
	sessionKeyUser  string = "userID"
	sessionKeyFlash        = "flash"
)

type App struct {
	UserService rsvp.UserService
	ErrorLog *log.Logger
	Session *sessions.Session
}

func (a App) home(w http.ResponseWriter, r *http.Request) {
	const op = "http.home"

	files := []string{
		"./pkg/ui/template/home.page.tmpl",
		"./pkg/ui/template/base.layout.tmpl",
		"./pkg/ui/template/footer.partial.tmpl",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (a App) showRegistrationForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showRegistrationForm"

	files := []string{
		"./pkg/ui/template/registration.page.tmpl",
		"./pkg/ui/template/base.layout.tmpl",
		"./pkg/ui/template/footer.partial.tmpl",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (a App) register(w http.ResponseWriter, r *http.Request) {
	const op = "http.register"

	r.FormValue("names")
	names := r.FormValue("names")
	email := r.FormValue("email")
	password := r.FormValue("password")

	u := rsvp.User{
		Names: names,
		Email: email,
	}

	_, err := a.UserService.CreateUser(&u, password)
	if err != nil {
		serverError(w, a.ErrorLog, err)
	}

	a.Session.Put(r, sessionKeyFlash, "Your account was successfully created. Please login")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a App) showLoginForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showLoginForm"

	files := []string{
		"./pkg/ui/template/login.page.tmpl",
		"./pkg/ui/template/base.layout.tmpl",
		"./pkg/ui/template/footer.partial.tmpl",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (a App) login(w http.ResponseWriter, r *http.Request) {
	const op = "http.login"

	r.FormValue("names")
	email := r.FormValue("email")
	password := r.FormValue("password")
	returnUrl := r.FormValue("return-url")

	if returnUrl == "" {
		returnUrl = "/"
	}

	// check if email and password match
	match, uid, err := a.UserService.EmailMatchPassword(email, password)
	if err != nil {
		fmt.Println(err)
		// todo:: handle error
		return
	} else if !match {
		fmt.Println(errors.New("unauthorized"))
		// handle error, user is unauthorized
		// todo:: re-render login page with errors
		return
	}

	a.Session.Put(r, sessionKeyUser, uid)

	// todo:: change to redirect to "returnUrl"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package http

import (
	"errors"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"rsvp/pkg/rsvp"
	"time"
)

var (
	sessionKeyUser  string = "userID"
	sessionKeyFlash        = "flash"
)

type templateData struct {
	User *rsvp.User
	Flash string
	CurrentYear int
	Data interface{}
}

type App struct {
	UserService rsvp.UserService
	ErrorLog *log.Logger
	Session *sessions.Session
	TemplateCache map[string]*template.Template
}

func (a App) home(w http.ResponseWriter, r *http.Request) {
	const op = "http.home"

	a.render(w, r, "home.page.tmpl", "hello")
}

func (a App) showRegistrationForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showRegistrationForm"

	a.render(w, r, "registration.page.tmpl", "Hello, world")
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

	a.render(w, r, "login.page.tmpl", "Hello, world")
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

func (a *App) logoutUser(w http.ResponseWriter, r *http.Request) {
	a.Session.Remove(r, sessionKeyUser)
	a.Session.Put(r, sessionKeyFlash, "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NewTemplateCache parses and caches templates in dir
func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	//the extension '.page.tmpl'. This essentially gives us a slice of all the
	//'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	return cache, nil
}

func (a App) render(w http.ResponseWriter, r *http.Request, name string, td interface{}) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper.
	ts, ok := a.TemplateCache[name]
	if !ok {
		// todo:: show server error
		//a.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, a.addDefaultData(td, r))
	if err != nil {
		// todo:: show server error
		//a.serverError(w, err)
		return
	}
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func (a *App) addDefaultData(data interface{}, r *http.Request) interface{} {
	u, _ := rsvp.UserFromContext(r.Context())

	td := &templateData{
		User: u,
		Flash: a.Session.PopString(r, "flash"),
		CurrentYear: time.Now().Year(),
		Data: data,
	}

	return td
}

package http

import (
	"errors"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"runtime/debug"
	"strconv"
	"time"
)

var (
	sessionKeyUser  string = "userID"
	sessionKeyFlash        = "flash"
)

type templateData struct {
	User *events.User
	Flash string
	CurrentYear int
	Data interface{}
}

type App struct {
	ErrorLog      *log.Logger
	Session       *sessions.Session
	TemplateCache map[string]*template.Template
	UserService   events.UserService
	EventService  events.EventService
}

func (a App) home(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "home.page.tmpl", nil)
}

func (a App) showRegistrationForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "registration.page.tmpl", nil)
}

func (a App) register(w http.ResponseWriter, r *http.Request) {
	const op = "http.register"

	r.FormValue("names")
	names := r.FormValue("names")
	email := r.FormValue("email")
	password := r.FormValue("password")

	u := events.User{
		Names: names,
		Email: email,
	}

	_, err := a.UserService.CreateUser(&u, password)
	if err != nil {
		a.serverError(w, r, err)
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
	a.Session.Put(r, sessionKeyFlash, "You were successfully logged in")

	// todo:: change to redirect to "returnUrl"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) logoutUser(w http.ResponseWriter, r *http.Request) {
	a.Session.Remove(r, sessionKeyUser)
	a.Session.Put(r, sessionKeyFlash, "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a App) showEvents(w http.ResponseWriter, r *http.Request) {
	const op = "http.showEvents"

	a.render(w, r, "events.page.tmpl", nil)
}

func (a App) showEventForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showEventForm"

	a.render(w, r, "event-form.page.tmpl", nil)
}

func (a App) createEvent(w http.ResponseWriter, r *http.Request) {
	const op = "http.createEvent"

	i, err := eventFromRequest(r)
	if err != nil {
		// todo:: handle
		fmt.Println(err)
		return
	}

	u, ok := events.UserFromContext(r.Context())
	if !ok {
		// todo:: handle -> server error
		fmt.Println(err)
		return
	}

	_, err = a.EventService.CreateEvent(i, u.ID)
	if err != nil {
		a.serverError(w, r, err)
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func eventFromRequest(r *http.Request) (*events.Event, error) {
	const op = "eventFromRequest"

	var id int
	var err error
	if idStr := r.FormValue("id"); idStr == "" {
		id = 0
	} else {
		id, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return nil, errors2.Wrap(err, op, "getting id form value")
		}
	}

	tz := r.FormValue("timezone")
	if tz == "" {
		return nil, errors2.Wrap(err, op, "getting timezone form value")
	}

	var startTime *time.Time
	if startTimeStr := r.FormValue("start_time"); startTimeStr == "" {
		startTime = nil
	} else {
		//t, err := time.Parse(time.RFC3339, startTimeStr)
		t, err := time.Parse("2006-01-02T15:04", startTimeStr)
		if err != nil {
			return nil, errors2.Wrap(err, op, "parsing form time")
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			return nil, errors2.Wrap(err, op, "parsing form time")
		}

		t = t.In(loc)
		startTime = &t
	}

	var endTime *time.Time
	if endTimeStr := r.FormValue("end_time"); endTimeStr == "" {
		endTime = nil
	} else {
		//t, err := time.Parse(time.RFC3339, endTimeStr)
		t, err := time.Parse("2006-01-02T15:04", endTimeStr)
		if err != nil {
			return nil, errors2.Wrap(err, op, "parsing form time")
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			return nil, errors2.Wrap(err, op, "parsing form time")
		}

		t = t.In(loc)
		endTime = &t
	}

	var isPublished bool
	isPublishedStr := r.FormValue("is_published")
	if isPublishedStr == "true" {
		isPublished = true
	} else if isPublishedStr == "false" ||isPublishedStr == "" {
		isPublished = false
	} else {
		return nil, errors2.Wrap(
			fmt.Errorf(
				"is_published form value could not be converted to boolean, got %v", isPublishedStr),
				op,
				"getting is_published form value",
		)
	}

	return &events.Event{
		ID:             id,
		Title:          r.FormValue("title"),
		Description:    r.FormValue("description"),
		IsVirtual:      false,
		Address:        r.FormValue("address"),
		Link:           r.FormValue("link"),
		NumberOfSeats:  0,
		StartTime:      startTime,
		EndTime:        endTime,
		WelcomeMessage: r.FormValue("welcome_message"),
		IsPublished:    isPublished,
	}, nil
}

func (a *App) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	a.render(w, r, "404.page.tmpl", nil)
}

func (a *App) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = a.ErrorLog.Output(2, trace)
	a.render(w, r, "server-error.page.tmpl", nil)
	w.WriteHeader(http.StatusInternalServerError)
}

func (a *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// NewTemplateCache parses and caches templates in dir
func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
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
	u, _ := events.UserFromContext(r.Context())

	td := &templateData{
		User: u,
		Flash: a.Session.PopString(r, "flash"),
		CurrentYear: time.Now().Year(),
		Data: data,
	}

	return td
}

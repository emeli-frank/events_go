package http

import (
	"errors"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/form"
	"events/pkg/validation"
	"fmt"
	"github.com/golangcollege/sessions"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	sessionKeyUser  string = "userID"
	sessionKeyFlash        = "flash"
)

const (
	defaultMultipartFormMaxMemory = 32 << 20 // 32 MB
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
	UploadDir	  string
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
		a.serverError(w, r, err)
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

	u, ok := events.UserFromContext(r.Context())
	if !ok {
		a.serverError(w, r, errors.New("couldn't get user from request context"))
		return
	}

	ee, err := a.EventService.Events(u.ID)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.render(w, r, "event-list.page.tmpl", struct {
		Events []events.Event
		ImageBaseURL string
	}{ee, a.UploadDir})
}

func (a App) showEvent(w http.ResponseWriter, r *http.Request) {
	const op = "http.showEvent"

	id, err := strconv.Atoi(mux.Vars(r)["eventID"])
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	u, ok := events.UserFromContext(r.Context())
	if !ok {
		a.serverError(w, r, errors.New("couldn't get user from request context"))
		return
	}

	e, err := a.EventService.Event(id)
	if err != nil {
		if _, ok := errors2.Unwrap(err).(*events.NotFound); ok {
			a.notFound(w, r)
			return
		}
		a.serverError(w, r, err)
		return
	} else if e.HostID != u.ID {
		a.notFound(w, r)
		return
	}

	a.render(w, r, "event-detail.page.tmpl", struct {
		Event *events.Event
		ImageBaseURL string
	}{e, a.UploadDir})
}

func (a App) showEventCreationForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showEventCreationForm"

	a.render(w, r, "create-event.page.tmpl", form.New(r.Form))
}

func (a App) createEvent(w http.ResponseWriter, r *http.Request) {
	const op = "http.createEvent"

	const maxFormSize = 2 << 19
	var hasFile bool

	err := r.ParseMultipartForm(defaultMultipartFormMaxMemory)
	if err != nil {
		a.Session.Put(r, sessionKeyFlash, "There was an error processing the form you submitted")
		fmt.Println(err)
		// todo:: show user a message
		a.render(w, r, "create-event-page.tmpl", nil)
		return
	}

	e, err := eventFromRequest(r)
	if err != nil {
		a.Session.Put(r, sessionKeyFlash, "There was an error processing the form you submitted")
		fmt.Println(err)
		// todo:: show user a message
		a.render(w, r, "create-event-page.tmpl", nil)
		return
	}

	file, handler, err := r.FormFile("cover_image")
	if err == http.ErrMissingFile {
		hasFile = false
	} else if err != nil {
		a.serverError(w, r, err)
		return
	} else {
		hasFile = true
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	f := form.New(r.PostForm)

	err = f.ValidateField("title", "", validation.Required, validation.Length(0, 64)).
		ValidateField("description", "", validation.Length(0, 512)).
		ValidateField("address", "", validation.Length(0, 128)).
		ValidateField("link", "", validation.Length(0, 128)).
		ValidateField("number_of_seats", "", validation.Length(2, 2 << 16)).
		ValidateField("start_time", "", events.StartTimeRule(false)).
		ValidateField("end_time", "", events.EndTimeRule(e.StartTime, false)).
		ValidateField("welcome_message", "", validation.Length(0, 256)).
		//ValidateField("cover_image", "Image is too large", validation.InRange(0, int(handler.Size))). // todo:: validate size but be mindful that handler may be nil, consider using a custom validator
		ValidateField("cover_image", "Image is too large", validation.RuleFunc(func(v interface{}) error { // todo:: find a better way to validate
			if !hasFile {
				return nil
			}

			size, ok := v.(int64)
			if !ok {
				return fmt.Errorf("expected int64, got %T\n", size)
			}

			if size > maxFormSize {
				return fmt.Errorf("Image size cannot be more than %d \n", size)
			}
			return nil
			}),
		).
		// todo:: validate scope[emails, ""]
		Error()
	//fmt.Printf(">>>>>>>>>>>>>>>>>>>>> type: %T, value: %p\n", err, err)
	if err != nil {
		if e, ok := err.(form.Error); ok {
			fmt.Println("form error: ", e.ErrorMessages()) // todo:: handler
			fmt.Println(f)
			a.render(w, r, "create-event.page.tmpl", f)
			return
		} else {
			a.serverError(w, r, err)
			return
		}
	}

	u, ok := events.UserFromContext(r.Context())
	if !ok {
		a.serverError(w, r, errors.New("could not get user from context"))
		return
	}

	var fileBytes []byte
	var ext string
	// read bytes from uploaded file
	if hasFile {
		// todo:: use something more secure. Look for a lib that can do this.
		ext, err = fileExt(handler.Header.Get("Content-Type"))
		if err != nil {
			a.Session.Put(r, sessionKeyFlash, "Unsupported file type uploaded")
			a.render(w, r, "create-event.page.tmpl", nil) // todo:: add data
			return
		}

		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			a.serverError(w, r, err)
			return
		}
	}

	_, err = a.EventService.CreateEvent(e, fileBytes, ext, u.ID)
	if err != nil {
		if _, ok := err.(validation.Error); ok {
			fmt.Println("validation error")
			// todo:: handle, pass back to form
		}
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func (a App) showEventEditForm(w http.ResponseWriter, r *http.Request) {
	const op = "http.showEventEditForm"

	a.render(w, r, "edit-event.page.tmpl", form.New(r.Form))
}

func (a App) editEvent(w http.ResponseWriter, r *http.Request) {
	const op = "http.editEvent"

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func (a App) publishEvent(w http.ResponseWriter, r *http.Request) {
	const op = "http.publishEvent"

	id, err := strconv.Atoi(mux.Vars(r)["eventID"])
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	u, ok := events.UserFromContext(r.Context())
	if !ok {
		a.serverError(w, r, errors.New("could not get user from request context"))
	}

	err = a.EventService.PublishEvent(id, u.ID)
	if err != nil {
		if _, ok := err.(validation.Error); ok {
			a.Session.Put(r, sessionKeyFlash,
				"Form could not be published. Please complete the form and try again")
			http.Redirect(w, r, "/events/edit/", http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func (a App) test(w http.ResponseWriter, r *http.Request) {
	const op = "http.test"

	a.serverError(w, r, errors.New("something very bad"))
	return
}

func eventFromRequest(r *http.Request) (*events.Event, error) {
	const op = "eventFromRequest"

	var err error
	var id int
	if idStr := r.PostForm.Get("id"); idStr == "" {
		id = 0
	} else {
		id, err = strconv.Atoi(r.PostForm.Get("id"))
		if err != nil {
			return nil, errors2.Wrap(err, op, "getting id form value")
		}
	}

	tz := r.PostForm.Get("timezone")
	if tz == "" {
		return nil, errors2.Wrap(errors.New("timezone was not set"), op, "getting timezone form value")
	}

	var startTime *time.Time
	if startTimeStr := r.PostForm.Get("start_time"); startTimeStr == "" {
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
	if endTimeStr := r.PostForm.Get("end_time"); endTimeStr == "" {
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
	isPublishedStr := r.PostForm.Get("is_published")
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
		Title:          r.PostForm.Get("title"),
		Description:    r.PostForm.Get("description"),
		Link:           r.PostForm.Get("link"),
		StartTime:      startTime,
		EndTime:        endTime,
		WelcomeMessage: r.PostForm.Get("welcome_message"),
		IsPublished:    isPublished,
	}, nil
}

func (a *App) notFound(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(fmt.Errorf("server error: the template %s does not exist", name))
		a.serverError(w, r, fmt.Errorf("server error: the template %s does not exist", name))
		//a.serverError(w, r, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, a.addDefaultData(td, r))
	if err != nil {
		a.serverError(w, r, err)
		fmt.Println(fmt.Errorf("server error: %v", err))
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

	return t.UTC().Format("02 Jan 2006, 15:04")
}

func (a *App) addDefaultData(data interface{}, r *http.Request) *templateData {
	u, _ := events.UserFromContext(r.Context())

	td := &templateData{
		User: u,
		Flash: a.Session.PopString(r, "flash"),
		CurrentYear: time.Now().Year(),
		Data: data,
	}

	return td
}

func fileExt(mime string) (string, error) {
	switch mime {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	default:
		return "", errors.New("unexpected mime type")
	}
}

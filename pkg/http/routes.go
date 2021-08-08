package http

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (a App) Routes(uploadDir string) http.Handler {
	standardMiddleWare := alice.New(/*h.recoverPanic ,*/ /*h.setReqCtxUser*/)
	//authOnlyMiddleWare := alice.New(/*s.checkJWT, */s.authenticateUser)
	dynamicMiddleware := alice.New(a.Session.Enable, a.addUserToSession)
	authenticatedOnly := alice.New(dynamicMiddleware.Then, a.authenticatedUser)

	r := mux.NewRouter()
	r.Handle("/", dynamicMiddleware.Then(http.HandlerFunc(a.home))).Methods("GET")
	r.Handle("/register", dynamicMiddleware.Then(http.HandlerFunc(a.showRegistrationForm))).Methods("GET")
	r.Handle("/register", dynamicMiddleware.Then(http.HandlerFunc(a.register))).Methods("POST")
	r.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(a.showLoginForm))).Methods("GET")
	r.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(a.login))).Methods("POST")
	r.Handle("/logout", dynamicMiddleware.Then(http.HandlerFunc(a.logoutUser))).Methods("GET")
	r.Handle("/events", authenticatedOnly.Then(http.HandlerFunc(a.showEvents))).Methods("GET")
	r.Handle("/events/{eventID:[0-9]+}", authenticatedOnly.Then(http.HandlerFunc(a.showEvent))).Methods("GET")
	r.Handle("/events/create", authenticatedOnly.Then(http.HandlerFunc(a.showEventCreationForm))).Methods("GET")
	r.Handle("/events/create", authenticatedOnly.Then(http.HandlerFunc(a.createEvent))).Methods("POST")
	r.Handle("/events/{eventID:[0-9]+}/edit", authenticatedOnly.Then(http.HandlerFunc(a.showEventEditForm))).Methods("GET")
	r.Handle("/events/{eventID:[0-9]+}/edit", authenticatedOnly.Then(http.HandlerFunc(a.updateEvent))).Methods("POST")
	r.Handle("/events/{eventID:[0-9]+}/delete", authenticatedOnly.Then(http.HandlerFunc(a.deleteEvent))).Methods("POST")

	r.NotFoundHandler = dynamicMiddleware.Then(http.HandlerFunc(a.notFound))

	// file route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/static"))))

	// upload route
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// test path
	r.Handle("/test", dynamicMiddleware.Then(http.HandlerFunc(a.test))).Methods("GET")

	return standardMiddleWare.Then(r)
}

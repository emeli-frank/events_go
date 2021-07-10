package http

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (a App) Routes() http.Handler {
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
	r.Handle("/events/create", authenticatedOnly.Then(http.HandlerFunc(a.showEventForm))).Methods("GET")
	r.Handle("/events/create", authenticatedOnly.Then(http.HandlerFunc(a.createEvent))).Methods("POST")

	r.NotFoundHandler = dynamicMiddleware.Then(http.HandlerFunc(a.notFoundHandler))

	// file route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/ui/static"))))

	return standardMiddleWare.Then(r)
}

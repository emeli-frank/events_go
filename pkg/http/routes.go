package http

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (a App) Routes() http.Handler {
	standardMiddleWare := alice.New(/*h.recoverPanic ,*/ /*h.setReqCtxUser*/)
	//authOnlyMiddleWare := alice.New(/*s.checkJWT, */s.authenticatedOnly)
	//dynamicMiddleware := alice.New(a.Session.Enable)

	r := mux.NewRouter()

	r.Handle("/", a.Session.Enable(http.HandlerFunc(a.home))).Methods("GET")
	r.Handle("/register", a.Session.Enable(http.HandlerFunc(a.showRegistrationForm))).Methods("GET")
	r.Handle("/register", a.Session.Enable(http.HandlerFunc(a.register))).Methods("POST")
	r.Handle("/login", a.Session.Enable(http.HandlerFunc(a.showLoginForm))).Methods("GET")
	r.Handle("/login", a.Session.Enable(http.HandlerFunc(a.login))).Methods("POST")

	// file route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/ui/static"))))

	return standardMiddleWare.Then(r)
}

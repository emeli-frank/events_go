package http

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (a App) Routes() http.Handler {
	standardMiddleWare := alice.New(/*h.recoverPanic ,*/ /*h.setReqCtxUser*/)
	//authOnlyMiddleWare := alice.New(/*s.checkJWT, */s.authenticatedOnly)
	dynamicMiddleware := alice.New(a.Session.Enable, a.authenticate)

	r := mux.NewRouter()
	r.Handle("/", dynamicMiddleware.Then(http.HandlerFunc(a.home))).Methods("GET")
	r.Handle("/register", dynamicMiddleware.Then(http.HandlerFunc(a.showRegistrationForm))).Methods("GET")
	r.Handle("/register", dynamicMiddleware.Then(http.HandlerFunc(a.register))).Methods("POST")
	r.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(a.showLoginForm))).Methods("GET")
	r.Handle("/login", dynamicMiddleware.Then(http.HandlerFunc(a.login))).Methods("POST")
	r.Handle("/logout", dynamicMiddleware.Then(http.HandlerFunc(a.logoutUser))).Methods("GET")

	// file route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/ui/static"))))

	return standardMiddleWare.Then(r)
}

package http

import (
	"fmt"
	"net/http"
	errors2 "events/pkg/errors"
	"events/pkg/events"
)

func (a App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				switch t := err.(type) {
				case error:
					//h.Response.serverError(w, t)
					panic(err) // todo:: remove and implement properly
				default:
					msg := fmt.Sprint("an unknown error:", t)
					_ = msg
					//h.Response.serverError(w, errors.New(msg))
					panic(err) // todo:: remove and implement properly
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a App) authenticatedUser(next http.Handler) http.Handler {
	const op = "app.authenticatedUser"

	f := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := events.UserFromContext(r.Context()); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

// addUserToSession adds user to session if their id is present in
// session and calls next.ServerHTTP or just calls next.ServerHTTP if
// their id is absent.
func (a *App) addUserToSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.Session.Exists(r, sessionKeyUser)
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.UserService.User(a.Session.GetInt(r, sessionKeyUser))
		if _, ok := errors2.Unwrap(err).(*events.NotFound); ok {
			a.Session.Remove(r, sessionKeyUser)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			// todo:: show server error
			fmt.Println(err)
			return
		}

		ctx := events.NewUserContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

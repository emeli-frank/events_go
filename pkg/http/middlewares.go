package http

import (
	"fmt"
	"net/http"
	"rsvp/pkg/rsvp"
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

/*func (h Http) authenticatedOnly(next http.Handler) http.Handler {
	const op = "server.authenticatedOnly"

	f := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := ecommerce.UserFromContext(r.Context()); !ok {
			h.Response.clientError(w, http.StatusUnauthorized, "")
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}*/

func (a *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := a.Session.Exists(r, sessionKeyUser)
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.UserService.User(a.Session.GetInt(r, sessionKeyUser))
		if /*err == models.ErrNoRecord*/ false { // todo:: check for not found error
			a.Session.Remove(r, sessionKeyUser)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			// todo:: show server error
			fmt.Println(err)
			return
		}

		ctx := rsvp.NewUserContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

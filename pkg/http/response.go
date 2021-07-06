package http

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func serverError(w http.ResponseWriter, errorLog *log.Logger, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = errorLog.Output(2, trace)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("A server error occurred"))
}

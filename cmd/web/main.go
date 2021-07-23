package main

import (
	http2 "events/pkg/http"
	"events/pkg/services"
	"events/pkg/storage"
	"events/pkg/storage/postgres"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", ":5000", "HTTP network address")
	sessionKey := flag.String("sessionKey", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session key")
	dsn := flag.String("dsn", "host=localhost port=5432 user=events password=password dbname=events sslmode=disable", "Postgresql database connection info")
	flag.Parse()

	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := storage.OpenDB("postgres", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte(*sessionKey))
	session.Lifetime = 12 * time.Hour
	//session.Lifetime = 5 * time.Second // todo:: remove

	// Initialize a new template cache
	templateCache, err := http2.NewTemplateCache("./pkg/static/template/")
	if err != nil {
		errorLog.Fatal(err)
	}

	psql := postgres.New(db)

	userRepo, err := postgres.NewUserStorage(psql)
	if err != nil {
		panic(err)
	}
	userService := services.NewUserService(userRepo)

	invitationRepo, err := postgres.NewEventStorage(psql)
	if err != nil {
		panic(err)
	}
	invitationService := services.NewEventService(invitationRepo)

	app := &http2.App{
		UserService:   userService,
		EventService:  invitationService,
		ErrorLog:      errorLog,
		Session:       session,
		TemplateCache: templateCache,
	}

	router := app.Routes()

	srv := &http.Server{
		Addr: *addr,
		Handler: router,
		ErrorLog: errorLog,
	}

	fmt.Printf("Starting server on: %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"rsvp/pkg/mock"
	"rsvp/pkg/storage"
	"time"
)

type application struct {
	DB                *sql.DB
	errorLog          *log.Logger
	infoLog           *log.Logger
	mocker 		 *mock.Mock
}

func main() {
	dsn := flag.String("dsn", "host=localhost port=5432 user=rsvp password=password dbname=rsvp sslmode=disable", "Postgresql database connection info")
	sqlScriptPath := flag.String("sql_script_path", "./pkg/storage/postgres/.db_setup/", "sql script path")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	sqlScriptPath2, err := filepath.Abs(*sqlScriptPath)
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := storage.OpenDB("postgres", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	mocker := &mock.Mock{
		DB:                db,
		W:                 os.Stdout,
	}

	app := application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		DB:               db,
		mocker: mocker,
	}

	fmt.Println("1. Create core data")
	fmt.Println("2. Create core and mock data")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	if choice == 0 {
		app.errorLog.Fatal("Please enter a valid choice")
	}

	switch choice {
	case 1: // create core data
		err = app.createCoreData(app, sqlScriptPath2)
		if err != nil {
			app.errorLog.Fatal(err)
		}
	case 2: // create core data and mock database
		err = app.createCoreData(app, sqlScriptPath2)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		err = app.createMockData(app, sqlScriptPath2)
		if err != nil {
			app.errorLog.Fatal(err)
		}
	default:
		app.errorLog.Fatal("Please enter a valid choice")
	}
}

func (a application) createCoreData(app application, sqlScriptPath string) error {
	fmt.Println("Creating tables...")
	t := time.Now()
	if err := app.createTables(sqlScriptPath); err != nil {
		return err
	}
	fmt.Printf("\tfinished creating table in %vs\n", time.Since(t))

	t = time.Now()
	fmt.Println("Creating core data...")
	err := storage.ExecScripts(a.DB, sqlScriptPath + "/data.sql")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\tfinished creating core data in %vs\n", time.Since(t))

	return nil
}

func (a application) createMockData(app application, sqlScriptPath string) error {
	t := time.Now()
	fmt.Println("Creating mock data...")
	err := a.mocker.SeedDB()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	fmt.Printf("\tfinished creating mock data in %vs\n", time.Since(t))

	return nil
}

func (a application) createTables(sqlScriptPath string) error {
	scriptPaths := []string{
		sqlScriptPath + "/teardown.sql",
		sqlScriptPath + "/tables.sql",
	}

	err := storage.ExecScripts(a.DB, scriptPaths...)
	if err != nil {
		panic(err)
	}

	return err
}

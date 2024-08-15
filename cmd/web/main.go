package main

import (
	"booking/internal/config"
	"booking/internal/driver"
	"booking/internal/handlers"
	"booking/internal/models"
	"booking/internal/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var InfoLog *log.Logger
var app config.AppConfig
var session *scs.SessionManager
var ErrorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)
	fmt.Println("Start email Listener ...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func run() (*driver.DB, error) {
	//register for models gonna put into session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	// change this to true when in production
	app.InProduction = false
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = InfoLog
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog
	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=07052002")
	log.Println("Database connected")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	return db, nil
}

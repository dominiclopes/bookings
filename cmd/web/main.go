package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dominiclopes/bookings/pkg/config"
	"github.com/dominiclopes/bookings/pkg/handlers"
	"github.com/dominiclopes/bookings/pkg/render"
)

const portNumber = ":8080"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("error creating the template cache, err: %v", err)
	}
	app.TemplateCache = tc

	render.NewAppConfig(&app)

	repo := handlers.NewRepository(&app)
	handlers.NewHandlers(repo)

	log.Printf("Starting application on port %s\n", portNumber)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting the server, err: %v\n", err)
	}
}

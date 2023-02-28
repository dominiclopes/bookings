package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

// writeToConsole add a log statement on every reqyest
func writeToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Hit the page. It's a %v request", r.Method)
		next.ServeHTTP(w, r)
	})
}

// noSurf adds CSRF protection to POST requests
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// sessionLoad loads and saves the session on every request
func sessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
